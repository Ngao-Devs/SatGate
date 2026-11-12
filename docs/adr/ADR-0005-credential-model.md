# ADR-0005: Credential and Authorization Model

- **Status:** Proposed
- **Date:** 2026-10-10

## Context

SatGate must verify that a protected email operation has satisfied its authorization requirements before allowing delivery.

The L402 protocol provides a mechanism for associating HTTP access with a Lightning payment through a macaroon and proof of payment. However, successful payment alone does not establish that a sender is legitimate or that a request is safe.

SatGate must evaluate how payment evidence and authorization rules should work together to protect against spam and repeated abuse.

## Decision

SatGate will investigate L402 credentials as a potential mechanism for verifying payment-linked authorization.

Regardless of the protocol selected, the authorization design must distinguish between:

- **Payment verification:** establishing whether the required payment has succeeded.
- **Credential verification:** establishing whether the presented credential is valid and its proof of payment is correct.
- **Request authorization:** determining whether the credential and applicable security rules permit the requested operation.

A protected email must not be sent unless the required authorization checks succeed.

## Authorization Requirements

The implementation must consider:

- whether the credential applies to the requested service and operation;
- whether the credential is valid and has not expired, where expiration applies;
- whether the associated payment has been successfully verified;
- whether usage restrictions or other applicable limits have been exceeded;
- whether additional anti-abuse checks are necessary.

The final credential format, validation procedure, and security policies will depend on the outcome of the L402 evaluation and the requirements of the initial email integration.

## Non-Goals

This ADR does not define:

- a new credential format;
- a final credential storage mechanism;
- a complete revocation system;
- a subscription or allowance implementation;
- the final L402 integration;
- the complete spam-detection strategy.

## Consequence

SatGate's authorization logic will not equate payment with trust. Payment evidence may satisfy an economic requirement, but the protected operation must still comply with its authorization rules and applicable anti-abuse controls.