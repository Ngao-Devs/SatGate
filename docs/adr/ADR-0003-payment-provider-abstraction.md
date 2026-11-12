# ADR-0003: Bitcoin Lightning Payment Architecture

- **Status:** Proposed
- **Date:** 2026-10-10

## Context

SatGate aims to protect email-sending operations from spam by introducing an economic barrier before a protected email can be sent.

Bitcoin's Lightning Network is a potential payment mechanism for this barrier. Lightning enables Bitcoin payments using satoshis (sats), making it a candidate for small-value payments where traditional payment methods may introduce disproportionate costs or friction.

SatGate must determine how Lightning payments can be incorporated into its security flow while ensuring that payment verification is reliable and that payment costs do not undermine the intended protection.

## Decision

SatGate will investigate the **Bitcoin Lightning Network as its initial candidate payment network**.

The evaluation will determine whether Lightning is suitable for collecting small payments associated with protected email submissions.

The design will account for:

- **Bitcoin (BTC):** the underlying monetary asset.
- **Satoshis (sats):** the unit used to express small Bitcoin-denominated payment amounts. One BTC equals 100,000,000 sats.
- **Lightning invoices:** payment requests that specify the amount and payment information required by the chosen Lightning implementation.
- **Payment verification:** confirmation that the required payment has succeeded before SatGate permits the protected email operation.

## Intended Payment Flow

```text
User Submits Contact Form
           ↓
      SatGate Check
           ↓
 Payment Requirement
           ↓
 Lightning Invoice (sats)
           ↓
   User Pays Invoice
           ↓
 Payment Verification
           ↓
   Email Authorization
           ↓
     Email Sent
```

This is the intended flow to investigate, not a claim that SatGate already implements it.

## Evaluation Criteria

The Lightning integration must be evaluated for:

- suitability for small-value payments;
- payment success and failure handling;
- reliable invoice settlement verification;
- transaction costs and operational expenses;
- integration complexity and service availability;
- the experience of legitimate users submitting contact forms;
- resistance to automated abuse by attackers willing to pay.

SatGate must also account for the possibility that payment costs, wallet requirements, or invoice-handling steps could discourage legitimate users.

## Non-Goals

This ADR does not yet select:

- a specific Lightning node or payment provider;
- a custody model or wallet implementation;
- the payment amount for each email;
- a final payment-to-access model;
- L402 as a mandatory protocol;
- a production deployment architecture.

## Consequence

Bitcoin Lightning becomes SatGate's initial payment technology to evaluate, with amounts expressed in satoshis. The final implementation will depend on whether it provides a practical, verifiable, and economically effective payment barrier for the anti-spam use case.