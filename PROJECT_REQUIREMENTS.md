# Referrences & Documents

1. https://docs.temporal.io/self-hosted-guide
2. https://docs.temporal.io/develop/go/
3. https://docs.temporal.io/evaluate/development-production-features/
4. https://github.com/temporalio/docker-compose
5. https://grafana.com/docs/k6/latest/set-up/install-k6/

# Targets

1. Go through Temporal development & production features.
2. Build a Saga Orchestrator Service for a microservices system by leveraging Temporal features.
3. Tech stacks:
   1. Backend with Go.
   2. Temporal.
   3. Transports:
      1. HTTP.
      2. Temporal signal, work queue.
   4. Testing:
      1. Load/performance test: K6.

# Product Requirements

## 1. Purpose

The goal of this project is to **learn Temporal’s workflow orchestration capabilities** through a fintech-style system. The system simulates end-to-end payment handling: request validation, debiting accounts, external payment gateway processing, fraud checks, and user notifications.

This project emphasizes **using all Temporal features** (workflows, activities, retries, timers, signals, queries, child workflows, versioning, task queues).

---

## 2. Business Requirements

### BR-1: Payment Request

- Users can initiate a payment by submitting a request with:
  - User ID
  - Amount
  - Payment method

### BR-2: Validation

- Validate user account and balance.
- If invalid, reject immediately.

### BR-3: Account Debit

- Debit user’s account before sending to external gateway.
- Roll back on failure.

### BR-4: External Payment Processing

- Attempt to process the payment via a (mock) external provider.
- If provider fails temporarily, retry with exponential backoff.

### BR-5: Fraud Check

- Every payment must pass an async fraud check service.
- Fraud check runs in parallel as a child workflow.

### BR-6: Gateway Callback

- The external provider sends async confirmation.
- System must accept **signals** for these callbacks.
- If no callback within 5 minutes, mark as failed.

### BR-7: Notifications

- Send user a notification (email/SMS) on success or failure.

### BR-8: Status Tracking

- Users can query the current status of their payment (pending, processing, success, fail).

### BR-9: Cancellation

- Users can signal the system to cancel an in-progress payment.

### BR-10: Workflow Evolution

- The system must support workflow **versioning** to allow safe upgrades in logic.

---

## 3. Functional Requirements

| Finish | Requirement                        | Temporal Feature        | Notes                         |
| ------ | ---------------------------------- | ----------------------- | ----------------------------- |
| ✅     | FR-1: Payment orchestration        | Workflow                | Central payment flow          |
| ✅     | FR-2: Validation, debit, credit    | Activities              | Simple DB ops                 |
| ✅     | FR-3: Retry external call          | Activity + Retry Policy | Retry w/ exponential backoff  |
| ✅     | FR-4: Fraud check                  | Child Workflow          | Runs async, result awaited    |
| ✅     | FR-5: Async gateway confirmation   | Signals                 | Handle external callbacks     |
| ✅     | FR-6: Timeout waiting for callback | Timers                  | Fail if >5 min                |
| ✅     | FR-7: Notification                 | Activity + Side Effect  | External email service        |
| ✅     | FR-8: Query payment status         | Query Handler           | Return current workflow state |
|        | FR-9: Cancel payment               | Signal Handler          | Abort gracefully              |
|        | FR-10: Workflow upgrades           | Versioning API          | Add new steps safely          |

---

## 4. Non-Functional Requirements

- Must run locally with Docker-based Temporal server.
- Written in a Temporal-supported SDK (Go, Java, Node, Python — pick one).
- Code should be modular: orchestrator, account service, fraud service, notification service.
- Mock external services (no real payment integration).

---

## 5. Learning Timeline

- **Week 1:** Setup Temporal + Hello World.
- **Week 2:** Implement Payment Orchestrator with basic validation + debit.
- **Week 3:** Add retries, signals, queries, timers.
- **Week 4:** Add child workflow (fraud check), notifications, cancellation, versioning.
- **Week 5:** Benchmark & testing.

# System design

### **Payment Orchestrator (Workflow)**

- Starts when a payment request comes in.
- Calls activities:
  - Validate request
  - Debit account
  - Call external payment gateway
- Uses Temporal features:
  - **Retries** for flaky external calls
  - **Timers** for waiting on async callbacks (e.g., fraud check)
  - **Signals** to receive gateway confirmation
  - **Queries** so clients can poll workflow state

### **Account Service (Activity workers)**

- Mock database for account balance.
- Activities: debit, credit.
- Show **idempotency** in activities.

### **Notification Service (Activity workers)**

- Sends user email/SMS when payment completes or fails.
- Demonstrates **side effects** (non-deterministic calls).

### **Fraud Check Service (Child Workflow)**

- Runs async fraud detection (mock delay).
- Returns pass/fail.
- Demonstrates **child workflows**.

### **Reporting Service (Query + Signals)**

- Query workflow state (pending, processing, success, fail).
- Send signal to cancel/abort a payment.

# Performance testing

## Load test

### Scenarios

- TODO: To be defined...

## Stress test

### Scenarios

- TODO: To be defined...
