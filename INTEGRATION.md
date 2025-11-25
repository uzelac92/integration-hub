#  **Integration Overview**

This document describes the full integration flow between the following components:

-   **Integration Hub** – receives debit/credit requests, validates signatures, enforces idempotency, and talks to operator
-   **Operator Mock (walletmock)** – simulated operator wallet with rate limits, random failures, and transaction log
-   **RGS (Real Game Server)** – receives webhook notifications
-   **Reconciler** – offline tool that compares hub transactions and operator transactions to detect mismatches

This document includes:

-   Architecture
-   Sequence diagrams
-   Idempotency guarantees
-   HMAC Signature scheme
-   Webhook Outbox & retry strategy
-   Reconciliation pipeline

# **Architecture**

```
+----------------------+          +---------------------+          +---------------------+
|    Integration Hub   | <------> |    Operator Mock    | <------> |         RGS         |
| (HMAC, idempotency)  |          | (wallet simulator)  |          |  (webhook receiver) |
+-----------+----------+          +----------+----------+          +----------+----------+
            |                                ^                               |
            |                                |                               |
            +-------- Webhook Outbox --------+------------- Webhook ---------+
            
+--------------------------------------------------------------------------+
|                                Reconciler                                |
| Fetch Hub Transactions + Fetch Operator Transactions + Diff + CSV Output |
+--------------------------------------------------------------------------+
```
# **HMAC Signature Scheme**

All incoming requests to Integration Hub must be signed with **HMAC-SHA256** using:
```
signature = HMAC_SHA256(secret, body + timestamp)
```
Headers used:
```
X-Timestamp: unix epoch seconds  
X-Signature: hex encoded HMAC
```

### **Verification Flow**

1.  Ensure both headers exist
2.  Convert timestamp to int
3.  Check timestamp skew ≤ 5 minutes
4.  Recompute HMAC
5.  Compare using hmac.Equal

# **Idempotency Flow**
Integration Hub enforces idempotency on all /wallet/* endpoints.

SQL table:
```
idempotency_keys(
    key TEXT UNIQUE,
    response JSONB,
    created_at TIMESTAMPTZ
)
```
Flow:
```
Client → Integration Hub → Idempotency Middleware
```
1.  Read Idempotency-Key header
2.  If key exists → return cached JSON immediately
3.  If not → run handler
4.  Save final response in DB
5.  Return response to client

# **Webhook Outbox Pattern**
Integration Hub uses an outbox table to reliably deliver webhook events to RGS.
```
webhook_outbox(
  id,
  event_id,
  payload,
  status = 'PENDING',
  attempt_count,
  next_attempt_at
)
```
### **Retry Loop (every 5 seconds)**

1.  SELECT all due events:
    status = PENDING AND next_attempt_at <= NOW()
2.  Try POSTing payload to RGS /webhook/hub
3.  If success → status = SUCCESS
4.  If fail →
    -   status = FAILED
    -   attempt_count += 1
    -   next_attempt_at = NOW() + exponential_backoff

# **Reconciliation Job**
The reconciler (cmd/reconciler) compares:
-   **Hub transactions** stored in PostgreSQL
-   **Operator transactions** collected by the Operator Mock

Operator exposes:
```
GET /v2/reconciliation
```
Reconciler performs:
1.  Fetch all hub transactions
2.  Fetch operator transactions
3.  Match by refId
4.  Detect mismatches:
    -   amount
    -   type
    -   operator balance
5.  Produce reconciliation_result.csv
6.  Exit:
    -   exit 1 if mismatches exist
    -   exit 0 if all matched

# **Sequence Diagram**

**Debit Flow (Hub → Operator → Webhook → RGS)**

![Sequence Diagram](https://www.plantuml.com/plantuml/png/RP71Zjem48RlVehHddQbZRBQOwHK2w4KH15LqrOzS9aaKs9XRCkUa4Q8Twy3TIbLZpFw_7_-xJVFwBWxA84r6mU5agHPsB2KjRIe6HPTTJTlB3aCxDto8L0mazuYrosv1q1_6-_HpnzAqI1ZXPvWDXNYweJatQZAuDEc_09fZqeHfmrLahVwSTdGmHecNG_9YePd-9wKSgUHEqTFOfJ7uzzGoH1Fi5XFYgt-C_wL024XdebjelVucbg50pcVdcHpJdl9RUEm5n6499fEu1cvyyzGeK9Tq-G7auEpYDcpQGPBxSanz5Irnko1ZcDLgTc0wRWoFwziWA-laC5SQbHmwGsEl3KXG8WyXuyufZ-Y7tJz17aknEs1esrEOGHBkf5w5wN-7p1yG6pxcqZRyCTXVRd83SxWzUK5Dgl_cTRezb40vTUYU5-MuYs8kuFw1G00)

# **Component Breakdown**

## **Integration Hub**

-   Validates signatures
-   Enforces idempotency
-   Sends debit/credit requests to operator
-   Stores all transactions in Postgres
-   Enqueues webhook events
-   Dispatches webhooks reliably (retry system)

## **Operator Mock**

-   Simulates real PSP
-   Has player balances
-   Stores all operator transactions
-   Supports rate limit (60 RPM)
-   Random internal 500 errors (10%)
-   Exposes /v2/reconciliation

## **RGS**

-   Receives incoming webhook notifications
-   Logs payload
-   Does not need business logic for this assignment

## **Reconciler**

-   Offline tool for mismatch detection
-   Produces CSV report
-   Returns a status code for CI pipelines