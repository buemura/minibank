# Feature Ideas

## Low-Hanging Fruit (gaps in existing functionality)

- [x] **Withdrawal endpoint** — `POST /api/v1/accounts/:id/withdrawals`
- [ ] **Wire up RabbitMQ** — Queue package and constants (`transaction.requested`) exist but async processing isn't connected in `svc-transaction`
- [x] **Use Redis caching** — Redis client is configured in the gateway but not actively used anywhere. Could cache account lookups, token validation, or balance queries
- [ ] **Account status management** — Domain supports `active`/`inactive`/`blocked` but there's no API to change status (close account, block/unblock)
- [ ] **User profile updates** — `svc-auth` has an `Update` repository method but no gRPC RPC or API endpoint to update name/phone
- [x] **Password change** — No endpoint to change password for authenticated users
- [x] **Fee transactions** — 5% fee automatically charged on withdrawal transactions

## Medium Effort

- [ ] **Rate limiting** — No API rate limiting on the gateway; could use Redis for distributed rate limiting
- [ ] **Transaction filtering** — History endpoint accepts `type_filter` in the proto but it's not wired through the gateway
- [ ] **Scheduled/recurring transfers** — Scheduled payments with a cron-based worker
- [ ] **2FA (TOTP)** — Two-factor authentication for login and sensitive operations
- [ ] **Webhook notifications** — Notify external systems on transaction events (completed, failed)
- [ ] **Admin API** — Admin endpoints for user management, account oversight, transaction monitoring

## Larger Features

- [ ] **Notification service** (`svc-notification`) — Email/SMS/push notifications for transactions, login alerts, low balance warnings. Natural consumer of RabbitMQ events
- [ ] **Multi-currency support** — Currency is hardcoded to BRL; could add exchange rates and cross-currency transfers
- [ ] **PIX integration** — QR code-based instant payments (fits the BRL context)
- [ ] **Audit logging service** — Immutable audit trail for all financial operations and auth events
- [ ] **KYC verification** — Document upload and identity verification flow
- [ ] **Spending analytics** — Monthly summaries, category-based spending, charts for the React frontend

## Infrastructure / Quality

- [ ] **Integration tests** — Current tests are all unit tests with mocks; no end-to-end gRPC or database tests
- [ ] **CI/CD pipeline** — No `.github/workflows` exist yet
- [ ] **Dockerize services** — Docker-compose only runs infra; services don't have Dockerfiles
- [ ] **Graceful error responses** — Gateway returns raw gRPC errors in some cases; could standardize error format
- [ ] **Request validation middleware** — Centralized input validation at the gateway level
- [ ] **Observability monitoring** — Re-enable Prometheus/Grafana (OpenTelemetry)
- [x] **Observability tracing** — Add distributed tracing (OpenTelemetry)
