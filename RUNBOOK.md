# **Overview**
This runbook covers:
-   How to diagnose common errors
-   Operator wallet failure flows
-   Idempotency & replay procedures
-   Webhook failures & replay
-   Reconciliation mismatches
-   Error code meanings
-   Manual recovery steps
-   Health & maintenance commands

Applies to:
-   **integration-hub**
-   **rgs**
-   **walletmock/operator**
-   **reconciler job**
-   **PostgreSQL (hubdb / rgsdb)**

# **Signature Errors**
Inbound /wallet/debit and /wallet/credit requests require:
-   X-Timestamp
-   X-Signature
-   Timestamp within ±5 minutes
-   Signature = HMAC(secret + body + timestamp)

### **Error: signature invalid**

**Symptoms:**
-   Hub logs:
    SIGNATURE MIDDLEWARE: signature invalid
-   Client receives **401 Unauthorized**

**Fix:**
1.  Check that body JSON **matches exactly** (spaces or newlines change signature).
2.  Verify timestamp is current:
```
echo $(date +%s)
```
3. Recompute signature:
```
sig = HMAC_SHA256(secret, body + timestamp)
```
4.  Re-send request.

# **Webhook Outbox Failures**
Integration-Hub stores outbound webhook jobs in:
```
webhook_outbox(
  id,
  event_id,
  payload,
  status, -- PENDING | FAILED | SUCCESS
  attempt_count,
  next_attempt_at
)
```
Dispatcher sends:
```
POST /webhook/hub  -->  RGS
```

### **When webhooks fail:**
Dispatcher log:
```
dispatcher: RGS status: 500
dispatcher: error sending webhook
```
Record moves from:
-   PENDING → FAILED with next retry scheduled

## **Manual Replay Procedure**
To retry a failed webhook:

### **1. Using SQL:**
```
UPDATE webhook_outbox
SET status = 'PENDING',
    next_attempt_at = NOW()
WHERE id = <id>;
```
Dispatcher will pick it up within the next cycle (5s).

### **2. Using Adminer**

-   Open webhook_outbox table
-   Locate row
-   Set status = 'PENDING'

# **Reconciliation Mismatches**
The reconciler produces CSV:
```
reconciliation_result.csv
```

### **Mismatch Row Format**
| refId | playerId | type | hubAmount | opAmount | hubStatus | opBalance |

# **RGS → Hub Webhook Failures**
RGS has its own webhook events table:
```
webhook_events(
  id,
  operator_id,
  event_type,
  payload,
  status,        -- pending | failed | completed
  retries,
  next_retry_at,
  error_message
)
```
### **Replay Procedure (RGS → Hub)**
```
POST /webhooks/retry/{id}
```
Or manually:
```
UPDATE webhook_events
SET status='pending', retries=0, next_retry_at=NOW()
WHERE id=<id>;
```

# **Common Error Codes (Hub & RGS)**
| **Error Message**              | **Meaning**                               | **Fix**                           |
| ------------------------------ | ----------------------------------------- | --------------------------------- |
| invalid JSON                   | malformed request                         | Fix request body                  |
| missing playerId               | req invalid                               | Add playerId                      |
| missing refId                  | idempotency key missing in wallet request | Add refId                         |
| operator error: 400            | operator rejected transaction             | Check operator reason             |
| operator error: 500            | operator unavailable                      | retry; dispatcher will auto retry |
| signature invalid              | HMAC mismatch                             | Recompute signature               |
| timestamp out of range         | clock drift                               | Fix client time                   |
| POSTGRES connection refused    | DB not started yet                        | Wait for container health         |
| no configuration file provided | reconciler missing config                 | Ensure config folder copied       |

# **Manual Recovery — Full Flow**

### **A) Replay Debit/Credit (Hub → Operator)**
```
curl -X POST /wallet/debit \
 -H "Idempotency-Key: <same>" \
 -d '{ same body }'
 ```

### **B) Replay RGS Webhook to Hub**
 ```
 curl -X POST /webhooks/retry/<id>
 ```

### **C) Replay Hub Webhook to RGS**
Set:
```
status = 'PENDING'
next_attempt_at = NOW()
```

## **D) Fix corrupted operator balance**
```
restart walletmock
```

# **Reconciler Usage**
Run reconciler:
```
docker compose run --rm reconciler
```
Outputs:
-   reconciliation_result.csv
-   Exit code reflects mismatch status.