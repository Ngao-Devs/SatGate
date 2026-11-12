# ADR-0002: Evaluate L402 for Spam Prevention

- **Status:** Proposed
- **Date:** 2026-10-10

## Context

SatGate aims to reduce spam by introducing an economic barrier before protected email operations can be completed.

The L402 protocol combines HTTP `402 Payment Required`, Lightning invoices, and cryptographic credentials that allow a client to prove payment when requesting a protected resource.

However, enabling automated clients to pay for services does not inherently prevent spam. Attackers may also automate payments if the cost is low enough.

SatGate must therefore determine whether L402 can support its security objectives and whether additional controls are necessary.

## Decision

SatGate will **research and evaluate L402 as a potential component of its anti-spam authorization flow**.

L402 will not be adopted as a core architectural dependency until its suitability has been assessed against SatGate's threat model.

The evaluation will examine:

- how a protected service requests payment;
- how successful payment is verified before an email is sent;
- how credentials are validated and associated with the protected operation;
- whether automated attackers could cheaply repeat the payment flow;
- what additional controls may be required to limit abuse.

## Evaluation Flow

```text
Contact Form Submission
          ↓
     SatGate Check
          ↓
 Does L402 fit the requirement?
          ↓
   Payment Challenge
          ↓
 Payment and Verification
          ↓
  Authorization Decision
          ↓
 Email Allowed or Rejected
```

This represents the flow being investigated, not a finalized implementation.

## Evaluation Criteria

L402 will be assessed according to whether it can:

- enforce a verifiable payment requirement;
- prevent email delivery before authorization succeeds;
- integrate with the intended contact-form protection;
- impose meaningful economic friction on repeated abuse;
- operate without introducing excessive friction for legitimate users.

The evaluation must distinguish between proving that a payment occurred and proving that a request is safe to fulfill. Payment alone does not establish that a sender is legitimate.

## Non-Goals

This ADR does not commit SatGate to:

- adopting L402 in production;
- enabling unrestricted automated purchasing;
- treating payment as sufficient proof that a request is legitimate;
- defining the final credential implementation;
- selecting a Lightning provider.

## Consequence

SatGate can evaluate an existing protocol without assuming that it solves the spam problem by itself. The eventual architecture will be selected based on evidence from the evaluation and the security requirements of the initial use case.