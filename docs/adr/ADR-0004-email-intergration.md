# ADR-0004: Initial Email Integration

- **Status:** Proposed
- **Date:** 2026-10-10

## Context

Contact forms are a common entry point for automated spam. Attackers can repeatedly submit forms to generate unwanted emails, consume resources, and disrupt legitimate communication.

SatGate aims to protect these operations by placing an authorization check before email delivery. Bitcoin Lightning is being evaluated as a possible way to introduce economic friction, while L402 remains under investigation as a potential part of the authorization flow.

## Decision

The initial SatGate use case will be **protecting website contact forms and their associated email-sending operations from spam and automated abuse**.

SatGate should integrate with an application's existing email-delivery mechanism rather than replace the email provider.

The protection layer must evaluate a submission before the protected email is sent. If the applicable authorization requirements are not satisfied, SatGate must prevent that operation from proceeding.

The specific payment and authorization mechanism will follow the decisions established through the Lightning and L402 evaluations.

## Intended Flow

```text
Visitor Submits Contact Form
             ↓
       Application
             ↓
          SatGate
             ↓
   Authorization Check
             ↓
   Required Conditions Met?
         /         \
       Yes          No
        ↓            ↓
  Email Provider   Reject or
        ↓          Challenge
    Email Sent
```

The final challenge and payment flow will depend on the selected implementation.

## Goals

- Protect contact forms from automated spam.
- Enforce authorization before protected email delivery.
- Integrate with existing email providers.
- Keep legitimate submissions practical.
- Provide a demonstrable initial use case for SatGate.
- Measure how effectively the integration reduces abusive submissions.

## Non-Goals

The initial integration will not attempt to:

- build a new email-delivery service;
- guarantee the elimination of all spam;
- assume payment alone proves a sender is legitimate;
- require every application to use the same email provider;
- finalize the payment model or L402 implementation.

## Consequence

Email protection becomes the first practical demonstration of SatGate's security model. The integration should keep the protection logic separate from email delivery so that future integrations can protect other operations without rebuilding the entire system.