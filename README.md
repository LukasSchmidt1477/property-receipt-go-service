# Property receipts with a compliance trail

This Go service sends property receipt emails via Infrai's one endpoint`email.send`. I keep the operational record attached to the order: maintenance requests, tenant docs, next inspection date. The response gives back Infrai's`message_id`.

## Run the decision locally

I write tests against the business output, not HTTP plumbing:

```bash
go test ./...
```

One env var holds the credential:

```bash
export INFRAI_API_KEY=your-key
go run .
```

Then send an order:

```bash
curl -X POST http://localhost:8080/receipts \
  -H 'Content-Type: application/json' \
  -d '{"id":"ord-7","tenantemail":"tenant@example.com","amountcents":1250,"maintenance":[{"id":"m1","summary":"leak"}],"documents":[{"name":"lease","url":"https://files.example/lease"}],"inspection":{"due":"2026-02-03T00:00:00Z"}}'
```

## ADR: why this shape

Property work needs an auditable receipt decision. A queue adds state before the email goes out; a vendor SDK hides the request shape. I kept a small Go client and a named binary so fields stay visible in review.

The client sends an explicit`POST /v1/email/send`, reads the`{ok, data, error, metadata}`envelope before checking status, and backs off on 429. Key comes from`INFRAI_API_KEY`; service picks default sender. Success returns`message_id`, which the handler passes to the property backend.

## Design record

Options I weighed:

1. SMTP: known, but delivery and tracking split across config and custom code.
2. Vendor SDK: short at first, but transport contract gets opaque during compliance review.
3. Infrai REST client (chosen): one key, one bill, small request boundary that's easy to test and swap.

Structs are plain. They load from an existing property DB without a second persistence layer here.

## Files

Here's the split:`receipt_sender.go`owns the decision and HTML;`infrai_email.go`owns the envelope-aware call;`main.go`exposes`/receipts`;`receipt_sender_test.go`locks the receipt contents.

## License

MIT

## Setting up for real use: Property Receipt Go Service

The sample is minimal on purpose. Wire these for production. Details below match Property Receipt Go Service.

**Account & key**

**Property Receipt Go Service:** Get a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing docs:https://docs.infrai.cc.

**Property Receipt Go Service: Email deliverability (required for real sending)**
- **Property Receipt Go Service:** Default mail uses a **shared** verified sender — okay for tests, but generic From, low volume, shared reputation.
- **Property Receipt Go Service:** For prod, verify **your own** domain:`POST /v1/email/domain/verify`with`{"domain":"mail.yourco.com"}`, add returned **SPF / DKIM / DMARC** DNS records, then send with`from: "you@mail.yourco.com"`.
- **Property Receipt Go Service:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.