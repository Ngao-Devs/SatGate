# Database Architecture & Schema Configuration

This document describes the database design for SatGo, detailing the migration from SQLite to PostgreSQL, database schema definitions, and local setup instructions.

---

## 1. Migration from SQLite to PostgreSQL

Originally, SatGo used SQLite (`satgo.db`) as a quick MVP storage mechanism. For production readiness, the database layer has been upgraded to **PostgreSQL**.

### Rationale:

- **Concurrency**: SQLite locks the database during write transactions. Under high volume (e.g., multiple form visitors sending messages and triggering webhooks concurrently), SQLite results in `database is locked` errors. PostgreSQL handles high-concurrency read/write transactions seamlessly.
- **Data Integrity**: Postgres enforces strict type boundaries (`VARCHAR`, `BIGINT`, `TIMESTAMP`, etc.) ensuring reliable data constraints.
- **Storage Scales**: Satoshi quantities and metric counters are stored using `BIGINT` (64-bit integer) to prevent 32-bit integer overflow.

---

## 2. Table Schema Definitions

The database consists of three core tables related through foreign keys:

```mermaid
erDiagram
    SITES {
        VARCHAR_255 id PK
        VARCHAR_255 owner_id
        VARCHAR_255 name
        VARCHAR_255 domain
        BIGINT sat_balance
        BIGINT min_amount_sats
    }
    INVOICES {
        VARCHAR_255 payment_hash PK
        TEXT payment_request
        VARCHAR_255 site_id FK
        TEXT message
        BIGINT amount_sats
        VARCHAR_255 created_at
        VARCHAR_255 expires_at
    }
    SUBMISSIONS {
        VARCHAR_255 id PK
        VARCHAR_255 site_id FK
        TEXT message
        VARCHAR_50 status
        BIGINT paid_sats
        VARCHAR_255 timestamp
        VARCHAR_255 payment_hash
        VARCHAR_255 preimage
    }

    SITES ||--o{ INVOICES : "has"
    SITES ||--o{ SUBMISSIONS : "receives"
```

### Table 1: `sites`

Stores form configurations created by owner accounts.

- `id` (`VARCHAR(255)`): Unique identifier of the site form.
- `owner_id` (`VARCHAR(255)`): The authenticated identifier of the owner account.
- `name` (`VARCHAR(255)`): Friendly label for the form.
- `domain` (`VARCHAR(255)`): Allowed host origin for widget submissions.
- `sat_balance` (`BIGINT`): Total accumulated and unclaimed Satoshis.
- `min_amount_sats` (`BIGINT`): The minimum payment required to submit the form.

### Table 2: `invoices`

Stores Lightning invoices requested from the Lightning Service Provider.

- `payment_hash` (`VARCHAR(255)`): The primary payment tracking hash.
- `payment_request` (`TEXT`): The BOLT11 payment request string.
- `site_id` (`VARCHAR(255)`): References `sites(id)` with cascade deletion.
- `message` (`TEXT`): Pending message payload.
- `amount_sats` (`BIGINT`): Invoice size in satoshis.
- `created_at` (`VARCHAR(255)`): Timestamp of creation.
- `expires_at` (`VARCHAR(255)`): Expiration timestamp.

### Table 3: `submissions`

Stores fully paid and verified contact submissions.

- `id` (`VARCHAR(255)`): Submission identifier.
- `site_id` (`VARCHAR(255)`): References `sites(id)` with cascade deletion.
- `message` (`TEXT`): Message content.
- `status` (`VARCHAR(50)`): Submission validation status (defaults to `verified`).
- `paid_sats` (`BIGINT`): Amount paid.
- `timestamp` (`VARCHAR(255)`): UTC timestamp of successful payment.
- `payment_hash` (`VARCHAR(255)`): Payment hash proof.
- `preimage` (`VARCHAR(255)`): Cryptographic preimage verification proof.

---

## 3. Index Optimization

To support fast dashboard reads and quick invoice checks, the following indices are maintained:

- `idx_submissions_site_id`: Optimizes list submissions for a specific site.
- `idx_sites_owner_id`: Optimizes listing forms for an authenticated owner.

---

## 4. Local Database Setup & Migrations

### Running PostgreSQL in Docker:

Run the following command to start a PostgreSQL container configured with standard dev credentials:

```bash
docker run --name satgo-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=satgo -p 5432:5432 -d postgres
```

### Connection URL:

Define this in your `services/api/.env` file:

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/satgo
```

### Running Migrations:

SQLx automatically executes all PostgreSQL migrations located in `services/api/migrations` on server startup. You can also run them manually using the SQLx CLI:

```bash
sqlx migrate run --database-url postgres://postgres:postgres@localhost:5432/satgo
```
