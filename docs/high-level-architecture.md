# High-Level Architecture

## 1. Overview

SatGate is a security service designed to protect applications from spam, automated submissions, and other forms of request abuse. Its initial focus is on contact forms and email-sending endpoints, where automated abuse can generate unwanted messages, consume resources, and create unnecessary costs.

SatGate introduces economic friction into these workflows to make repeated abuse more expensive and less attractive. It is not intended to replace conventional security controls. Instead, it combines payment requirements with request validation and other anti-abuse measures.

The backend is primarily written in **Go**. The payment integration currently uses the **Blink API**, which has been integrated and is working during the testing stage.

This document describes the system's architectural direction and the currently confirmed implementation. Detailed design decisions are maintained separately in the Architecture Decision Records (ADRs).

## 2. Architectural Goals

SatGate's architecture should support the following goals:

- **Abuse prevention:** Make automated spam and repeated abusive submissions harder to carry out at scale.
- **Economic friction:** Introduce a cost to selected requests where appropriate.
- **Reliable payment handling:** Ensure payment status is checked before a protected operation is allowed to proceed.
- **Application integration:** Allow existing applications to integrate SatGate without having to redesign their entire request-handling workflows.
- **Security by design:** Avoid trusting incoming requests, payment claims, or client-side state without appropriate verification.
- **Maintainability:** Keep responsibilities clearly separated so that security, payment, and application-integration logic can evolve independently.
- **Extensibility:** Leave room for additional payment mechanisms, credential models, and protected use cases as the project develops.

## 3. Main Components

### 3.1 Go Backend

The Go backend is the core of SatGate. It coordinates the application's server-side processing and provides the foundation for integrating payment functionality and anti-abuse controls.

Its responsibilities include handling incoming requests, enforcing the applicable rules, coordinating payment-related operations, and determining whether a request can proceed.

As the system evolves, these responsibilities should remain clearly separated within the backend rather than becoming tightly coupled to a single payment provider or application integration.

### 3.2 Blink API Integration

SatGate currently integrates with the Blink API for invoice creation and payment processing during testing.

This integration provides the payment functionality needed to test SatGate's economic-friction approach against real payment flows. It is an existing, working part of the service—not merely a proposed future integration.

The integration must support reliable handling of payment outcomes. In particular, SatGate must distinguish between an invoice being created, a payment being completed, and a protected request being authorized to proceed.

Payment verification and request authorization are related but distinct responsibilities. A successful payment should only satisfy the applicable payment requirement; it should not automatically bypass other security checks.

### 3.3 Protected Applications and Endpoints

SatGate is intended to protect application features that accept user submissions or trigger potentially abusable actions.

Initial use cases include:

- Contact forms that receive messages from visitors.
- Email-sending endpoints that could otherwise be abused to generate unwanted messages.
- Other request-based workflows where automated abuse creates operational or financial costs.

The goal is to provide a protection layer that applications can integrate into their existing workflows. The exact integration mechanism and supported endpoint contracts should be documented against the current implementation.

### 3.4 Security and Anti-Abuse Controls

Economic friction is only one part of SatGate's protection model. The system must also account for requests that are malicious, invalid, duplicated, or otherwise suspicious.

Security controls should be designed to address risks such as repeated submissions, forged payment claims, request manipulation, and attempts to bypass payment requirements.

Relevant controls include request validation, rate limiting, trustworthy payment-status verification, and appropriate handling of repeated or invalid operations. The precise controls and their implementation status should be tracked separately from the high-level architecture.

## 4. High-Level Request and Payment Flow

The following describes the intended interaction between a protected application, SatGate, and Blink. Exact routes and internal processing steps should follow the current Go implementation.

1. **Request initiation:** A user initiates an action on an application protected by SatGate.
2. **Request evaluation:** SatGate evaluates the request against the applicable validation, security, and payment requirements.
3. **Invoice creation:** If payment is required, SatGate uses the Blink API to create an invoice.
4. **Payment:** The user completes the payment through the supported payment flow.
5. **Payment verification:** SatGate determines whether the required payment has actually succeeded, rather than relying solely on information supplied by the client.
6. **Authorization:** SatGate evaluates whether the request satisfies the applicable payment and security requirements.
7. **Protected operation:** The request proceeds only when the relevant requirements have been satisfied.

This flow is a conceptual overview, not a claim that every step is implemented through a separate service or endpoint.

## 5. Security Principles

SatGate's architecture should follow several core principles:

- **Never trust client-side payment claims.** Payment status must be established through a trusted verification mechanism.
- **Separate payment from authorization.** Payment does not automatically make a request safe or legitimate.
- **Prevent replay and duplication.** Repeated payment notifications or request attempts should not result in unintended duplicate processing.
- **Validate incoming requests.** Inputs and request context must be checked before protected operations are performed.
- **Protect credentials.** Payment API credentials and other secrets must remain on the server and must not be exposed to client applications.
- **Fail safely.** If payment status cannot be reliably established, SatGate must not treat the payment requirement as satisfied.

These principles guide implementation and review; the corresponding ADRs should define specific decisions where further detail is needed.

## 6. Current Implementation Status

| Component | Current status |
|---|---|
| Go backend | Primary backend implementation |
| Blink API integration | Working during the testing stage |
| Invoice creation and payment flow | Available for testing |
| Anti-abuse controls | Part of the system's protection model; individual implementations should be verified against the codebase |
| Production readiness | Requires validation of security, reliability, configuration, and operational requirements |

The use of Blink during testing does not, by itself, determine whether it will remain the production payment provider. That decision should be documented separately if and when it is made.

## 7. Architectural Decisions and Future Development

SatGate's detailed technical decisions are documented in its ADRs, including the economic approach to spam prevention, evaluation of L402, Lightning payment architecture, email integration, credential and authorization models, and subscriptions and allowances.

These documents should distinguish between decisions already implemented, proposals under evaluation, and work that remains outstanding.

Future development should prioritize strengthening the existing working payment flow, verifying that payment and authorization are handled safely, and making SatGate straightforward to integrate into protected applications. Additional infrastructure or architectural complexity should be introduced only when justified by concrete requirements.

## 8. Summary

SatGate combines a Go-based backend, a working Blink API payment integration in the testing stage, and a broader anti-abuse model intended to protect application endpoints from spam and automated misuse.

Its central architectural principle is that **payment introduces economic friction, while security controls determine whether a request should be trusted and allowed to proceed**. Keeping those responsibilities distinct will help SatGate evolve without making payment the sole defense against abuse.