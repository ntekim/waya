# WAYA (WayaGrid) ⚡️

> **The Orchestration Layer for African Finance.**
> Built for the Afriex Hackathon 2025.

## 🚀 The Problem
Afriex has built incredible rails for Peer-to-Peer (P2P) money transfer. But for businesses, NGOs, and platforms, moving money is still manual. 
- How does a US startup pay 50 Nigerian engineers instantly? 
- How does an NGO disperse funds to 1,000 recipients in Ghana?
- Doing this one by one on a mobile app is impossible.

## 💡 The Solution: Waya
**Waya** is a high-performance B2B Orchestration Infrastructure API. It plugs into the Afriex ecosystem to enable **Bulk Payouts**, **Smart Salary Splitting**, and **Natural Language Transactions**.

We provide the "Grid" that allows other platforms (HR Software, Gig Marketplaces) to plug into Afriex.

## 🛠 Tech Stack
- **Language:** Golang (1.23)
- **Architecture:** Hexagonal (Ports & Adapters)
- **Database:** SQLite (Dev) / Postgres (Prod) via **SQLC** (Type-safe SQL)
- **Framework:** Echo v4
- **Documentation:** Swagger / OpenAPI
- **Frontend:** Next.js + Chakra UI (Dashboard)

## ⚡️ Key Features
1.  **High-Concurrency Engine:** Uses Go Goroutines to process 500+ payouts/second.
2.  **Type-Safe Database:** Leveraging SQLC for compiled SQL queries.
3.  **Embeddable SDK:** A drop-in solution for any platform to start sending money to Africa.

## 🏃‍♂️ Quick Start

### Prerequisites
- Go 1.22+
- Make
- sqlc (`go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`)

### Installation
1. Clone the repo
2. Setup environment
   ```bash
   cp .env.example .env