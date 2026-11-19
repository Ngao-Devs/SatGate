# SatGo

SatGo is an embeddable, Lightning-powered contact form that leverages micro-payments as economic proof of human intent.

## Core Value Proposition

SatGo replaces friction-heavy, AI-vulnerable CAPTCHAs with a micro-tax on communication. By requiring a tiny Lightning payment (e.g., 1 to 10 satoshis), automated spam campaigns become cost-prohibitive, while legitimate human interactions remain unhindered. Website owners eliminate data clutter while capturing micro-revenue from high-intent messages.

## System Architecture

The workspace consists of two decoupled components engineered for scale and speed:

- **Backend (`services/api`)**: High-concurrency Rust Axum server utilizing Tokio and SQLx for lightning-fast message state management and webhook verification.
- **Frontend (`apps/web`)**: Next.js App Router dashboard and client widget featuring seamless WebLN integration and responsive fallback QR mechanics.

## Quickstart Setup

### Prerequisites

Ensure you have `pnpm` and the Rust toolchain installed.

### Installation

Clone the repository and pull the frontend dependencies:

```bash
pnpm install
```

### Execution

1. Initialize the Rust API server:

```bash
pnpm api:dev
```

2. In a separate terminal session, fire up the Next.js development client:

```bash
pnpm dev
```

### Target Environments

- **Marketing Landing Page**: http://localhost:3000
- **Management Dashboard**: http://localhost:3000/dashboard
- **Live Product Demo**: http://localhost:3000/demo
- **API Health Interface**: http://localhost:8080/health

## Development Configuration

Initialize localized environment files before custom production deployments:

```bash
cp services/api/.env.example services/api/.env
cp apps/web/.env.example apps/web/.env.local
```
