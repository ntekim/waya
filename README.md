Here is the updated, detailed `README.md`. This version is structured to serve both as a pitch deck and a comprehensive technical guide for any developer looking at your code.

### 📄 Updated `README.md`

```markdown
# WAYA (WayaGrid) ⚡️

> **The High-Performance Orchestration Layer for African Cross-Border Payouts.**
> Built for the Afriex Hackathon 2025.

---

## 🚀 Waya's Vision: One API to Pay Africa

The Afriex Business API requires a complex, three-step chain to complete a single payment: (1) Create Customer → (2) Create Payment Method → (3) Create Transaction.

**Waya** is the solution. It transforms this tedious, sequential workflow into a **single, concurrent API call**. It is designed to process mass payments (like payroll or vendor disbursements) across multiple African countries instantly.

### 💡 Core Value Proposition (For Developers)

Waya enables seamless integration of Afriex payouts by replacing the complex sequential logic with a powerful, parallel processing engine.

| Feature | Afriex Direct (Multi-Call) | Waya Orchestrator (Single-Call) |
| :--- | :--- | :--- |
| **Workflow** | 3 separate sequential API calls per recipient. | **1** Asynchronous API call handles all steps concurrently. |
| **Concurrency** | Limited by single thread loop. | Go **Goroutines** process payouts in parallel for instant volume. |
| **Status Model** | Must implement complex polling logic per recipient. | Single **Batch ID** for real-time status check on entire manifest. |
| **Scalability** | Designed for simple transactions. | Designed for B2B Payroll/Vendor disbursement volume. |

---

## 💻 Tech Stack & Architecture

Waya is built on a high-concurrency, clean-architecture stack suitable for enterprise finance platforms.

*   **Language:** **Golang (Go)** - Chosen for its exceptional performance, concurrency model (Goroutines), and speed.
*   **Web Framework:** **Echo v4** - A minimal, high-performance web framework for the API endpoints.
*   **Architecture:** **Hexagonal (Ports & Adapters)** - Ensures the core logic (the Orchestrator) is independent of external dependencies (Afriex API, SQLite/Postgres DB).
*   **Database Tooling:** **SQLC** - Generates type-safe Go code from raw SQL, eliminating common runtime database errors and maximizing performance.
*   **Frontend Demo:** **Next.js + Chakra UI** - Provides a simple, clean interface to visually demonstrate the API's performance.

## ⚙️ Quick Start (Developer Guide)

### Prerequisites

You must have **Go (1.22+)** and **Make** installed.

### 1. Project Setup & Config

```bash
# 1. Clone the repository
git clone [YOUR-REPO-URL] waya
cd waya

# 2. Setup Dependencies
go mod tidy

# 3. Rename environment file and set keys
cp app.env .env # Create the .env file
# **IMPORTANT:** Edit .env and paste your Afriex API Key
```

### 2. Generate Database Code

Waya uses SQLC, which requires generating Go code from your `.sql` queries.

```bash
# 1. Regenerate SQLC code and initial Swagger docs
make setup 
# This runs: tidy -> sqlc -> swagger

# 2. Remove old incompatible DB file (if switching schemas)
rm waya.db
```

### 3. Run the Backend Orchestrator

```bash
make run
# Server will start on http://localhost:8080/ (or your PORT in .env)
```

---

## 🔌 API Endpoints (The Waya Contract)

All endpoints are hosted under the base path `/api/v1`.

### 1. Bulk Payout Orchestration

This is the primary endpoint that triggers the concurrent Afriex payment flow.

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| **POST** | `/payouts` | Accepts a JSON batch of payments, saves to DB, and starts the concurrent 3-step Afriex process in a background Goroutine. |

### 2. Batch Status Check

Allows the client to poll for the real-time status of the payouts in the batch.

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| **GET** | `/payouts/{batch_id}` | Retrieves the aggregated status and all individual payout records for a given batch. |

### 3. Live Documentation

Once the server is running, visit the auto-generated Swagger page:
`http://localhost:8080/swagger/index.html`

---

## 🚀 Frontend Demo Instructions

The companion Next.js application visualizes the speed of the Orchestrator.

```bash
# In a separate terminal (outside the Go project)
cd frontend

# Install dependencies (if not done)
npm install

# Run the frontend application
npm run dev
# Open http://localhost:3000 in your browser
```
Use the frontend to paste your JSON payload and watch the Go terminal logs to see the concurrency in action!

---
*Built with ❤️ by [Jotham]*
```
