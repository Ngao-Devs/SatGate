# ADR-0001: SatGate — Economic Protection Against Spam

- **Status:** Proposed
- **Date:** 2026-10-10

## Context

Web applications expose resources such as contact forms and email-sending endpoints that can be abused by automated bots to generate large volumes of unwanted submissions.

Traditional defenses, including rate limiting and CAPTCHA, can help reduce abuse but may not fully address the cost of large-scale automated attacks.

SatGate explores using an economic barrier to make abusive submissions more expensive, helping protect email and other digital services from spam.

## Decision

SatGate will be developed as a **security service that uses economic friction to help prevent spam and automated abuse**.

The initial use case will focus on protecting website contact forms and email-sending operations. Before a protected email is sent, the sender must satisfy the applicable authorization requirements, which may include a payment.

The Bitcoin Lightning Network will be investigated as a potential payment mechanism. The L402 protocol will also be researched to determine whether it can support SatGate's security requirements.

Neither technology is considered a confirmed architectural choice until the research establishes its suitability.

## Initial Use Case

SatGate will explore protecting a contact form so that automated attackers cannot freely generate large volumes of email submissions.

The intended flow is:

```text
Contact Form Submission
          ↓
     SatGate Check
          ↓
Authorization Requirement
          ↓
 Verification / Payment
          ↓
    Email Permitted
```

The exact payment and verification flow will be determined through subsequent architectural decisions.

## Goals

- Reduce spam and automated abuse of protected email endpoints.
- Investigate economic barriers as a security mechanism.
- Evaluate Lightning payments and L402 against the threat model.
- Keep legitimate submissions practical for real users.
- Design the protection so it can eventually support other digital services.

## Non-Goals

The initial version will not attempt to:

- guarantee the elimination of all spam;
- build a new cryptocurrency or payment network;
- assume payment alone prevents automated abuse;
- commit to L402 before evaluating it;
- build a complete email delivery provider.

## Success Criteria

The initial implementation should demonstrate that SatGate can protect an email-sending operation by enforcing its configured authorization requirements before allowing the email to be sent.

Testing should include legitimate submissions and automated abuse scenarios to evaluate both security effectiveness and user friction.

## Consequence

SatGate's architecture and technology choices will be guided by its primary security objective: **making automated spam more difficult and costly without unnecessarily obstructing legitimate users**.