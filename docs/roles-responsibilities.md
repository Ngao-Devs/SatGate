# Roles and Responsibilities

## 1. Purpose

This document defines the technical roles required to develop and maintain SatGate, the responsibilities associated with each role, and the recommended number of contributors needed.

SatGate is a security service designed to reduce spam and automated abuse through economic friction and appropriate security controls. Its backend is primarily written in Go, and its current testing-stage payment integration uses the Blink API.

These roles support the implementation of SatGate's Architecture Decision Records (ADRs). The staffing numbers are recommended targets rather than strict limits. Contributors may support multiple roles depending on their skills and the project's needs.

## 2. Technical Roles

### 2.1 Architecture and System Design

**Recommended contributors: 1–2**

**Responsibilities:**
- Maintain a consistent, maintainable system architecture.
- Ensure implementations align with approved ADRs.
- Identify architectural dependencies and potential design conflicts.
- Evaluate proposed architectural changes and document significant decisions.
- Help keep integrations between backend, payment, security, and application components consistent.

**Related ADRs:** ADR-0001, ADR-0002, ADR-0003, ADR-0005, ADR-0006.

### 2.2 Backend Engineering

**Recommended contributors: 3**

**Responsibilities:**
- Develop and maintain SatGate's Go backend.
- Implement API functionality and core request-processing logic.
- Integrate payment, authorization, and application-integration components.
- Implement reliable error handling and appropriate data validation.
- Write and maintain backend tests.
- Keep backend code modular, readable, and maintainable.

**Related ADRs:** ADR-0001, ADR-0003, ADR-0004, ADR-0005, ADR-0006.

### 2.3 Security and Anti-Abuse Engineering

**Recommended contributors: 2**

**Responsibilities:**
- Identify threats and potential abuse scenarios affecting SatGate.
- Design and implement appropriate anti-spam and request-validation controls.
- Evaluate rate limiting, replay prevention, and other abuse-mitigation techniques.
- Review security-sensitive changes and identify bypass opportunities.
- Ensure payment requirements complement, rather than replace, security controls.

**Related ADRs:** ADR-0001, ADR-0002, ADR-0005.

### 2.4 Payment Systems Engineering

**Recommended contributors: 2**

**Responsibilities:**
- Maintain and improve the existing Blink API payment integration.
- Work on invoice creation and reliable payment-status verification.
- Handle payment failures, duplicate notifications, and unexpected payment states.
- Ensure payment outcomes are handled safely and consistently.
- Evaluate payment architecture changes and future payment-provider requirements.

**Related ADRs:** ADR-0002, ADR-0003, ADR-0006.

**Current implementation:** Blink API is integrated and working during the testing stage. Any proposed replacement or additional provider should be evaluated separately.

### 2.5 Protocol and Credential Engineering

**Recommended contributors: 2**

**Responsibilities:**
- Research and evaluate protocols that could strengthen SatGate's authorization model.
- Evaluate L402 and its suitability for SatGate's anti-abuse use cases.
- Define how credentials, proofs, and authorization requirements should be handled.
- Consider credential validation, expiry, replay resistance, and misuse.
- Ensure protocol decisions fit the existing payment and security architecture.

**Related ADRs:** ADR-0002, ADR-0005, ADR-0006.

L402 remains an evaluation candidate unless a subsequent architectural decision formally adopts it.

### 2.6 Email and Application Integration Engineering

**Recommended contributors: 2**

**Responsibilities:**
- Develop integrations between SatGate and protected application workflows.
- Support the initial focus on contact forms and email-sending endpoints.
- Ensure integration flows handle validation, payment requirements, and failures appropriately.
- Make integrations understandable and practical for application developers.
- Document integration requirements and test expected behavior.

**Related ADRs:** ADR-0001, ADR-0004, ADR-0005, ADR-0006.

### 2.7 QA and Reliability Engineering

**Recommended contributors: 1–2**

**Responsibilities:**
- Develop and maintain automated tests for critical system behavior.
- Test invoice and payment flows, including failures and duplicate events.
- Verify authorization and anti-abuse behavior.
- Identify regressions and reliability issues.
- Help ensure changes meet their acceptance criteria before release.

**Related ADRs:** ADR-0001 through ADR-0006, where applicable.

### 2.8 DevOps and Operations Engineering

**Recommended contributors: 1**

**Responsibilities:**
- Support development, testing, and deployment environments.
- Maintain CI/CD pipelines and automated quality checks.
- Manage application configuration and deployment secrets securely.
- Support logging, monitoring, and operational troubleshooting.
- Help prepare SatGate for reliable production deployment.

**Related ADRs:** ADR-0001 through ADR-0006, where operational requirements apply.

## 3. Staffing Guidelines

The recommended numbers indicate the desired capacity for each technical area, not the total number of unique people required by the project.

- Contributors may participate in more than one role when their skills and availability allow it.
- Security-sensitive and payment-related changes should receive appropriate review, preferably from another contributor.
- The team should prioritize filling critical responsibilities over creating separate teams for every role.
- Staffing can evolve as SatGate's implementation, complexity, and contributor base grow.

## 4. Shared Responsibility

SatGate is a shared project. These roles clarify technical responsibilities and help contributors understand where their work fits; they do not establish exclusive ownership of the project.

Contributors are encouraged to collaborate across roles, share knowledge, review one another's work, and raise architectural or security concerns when they arise.

The Architecture Decision Records define the relevant technical decisions. Contributors should consult the applicable ADRs when implementing features or proposing changes to the system.