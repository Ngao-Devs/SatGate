# ADR-0006: Subscription and Allowance Model

- **Status:** Proposed
- **Date:** 2026-10-10

## Context

SatGate aims to reduce spam by introducing an economic barrier before protected email operations can proceed.

Different payment models offer different trade-offs. Per-submission payments may make repeated abuse more expensive, while allowances and subscriptions can simplify legitimate repeated usage but may also provide attackers with access to send multiple unwanted emails.

SatGate must evaluate these models against its security objectives before implementing them.

## Decision

SatGate will evaluate subscriptions and usage allowances as **optional access models, not default features**.

The evaluation will consider whether each model can preserve meaningful economic friction while allowing legitimate users to submit messages conveniently.

Potential models include:

- **Per-submission payment:** each protected email requires a separate payment or equivalent authorization.
- **Usage allowance:** a payment permits a defined number of submissions.
- **Time-based access:** authorization remains valid for a defined period.
- **Subscription:** recurring payments maintain access under an agreed plan.

No subscription or allowance model is selected by this ADR.

## Security Considerations

Before implementing reusable access, SatGate must determine:

- how many email operations a payment can authorize;
- whether an allowance can be exhausted or abused through automated requests;
- how usage limits are enforced under concurrent requests;
- how expired access and failed renewals are handled;
- whether additional rate limits or anti-abuse checks are required.

A valid subscription or allowance must not automatically exempt a client from applicable security controls.

## Non-Goals

This ADR does not define:

- pricing or subscription tiers;
- a fixed allowance per payment;
- recurring billing implementation;
- a usage-tracking database schema;
- a final decision between per-submission and reusable access.

## Consequence

SatGate will not introduce subscriptions or reusable allowances solely for convenience. Their value must be demonstrated against the primary objective: reducing automated spam without imposing unnecessary friction on legitimate users.