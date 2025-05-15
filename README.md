# REST API for CRUD operations using Go

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Overview
A scalable REST API for managing users, built with Go. It supports CRUD operations, email and mobile number uniqueness validation. The API uses PostgreSQL for persistent storage.

### Key Features
- **User CRUD**: Create, read, update, and delete users with validated email uniqueness.
- **Tech Stack**: Go, Chi, Pgx (PostgreSQL).

## Prerequisites
- Go 1.20+
- PostgreSQL 13+
- Docker (optional, for running services)

## Setup

### 1. Clone the Repository
```bash
git clone https://github.com/upekZ/rest-api-go.git
cd rest-api-go
```

### 2. Install Dependencies
```bash
go mod tidy
go get github.com/go-chi/chi/v5
go get github.com/jackc/pgx/v5
```

### 3. Configure Environment
Create a `.env` file:
```env
DATABASE_URL=postgres://user:password@localhost:5432/mydb
```

### 4. Set Up PostgreSQL
Run PostgreSQL (e.g., via Docker):
```bash
docker run -d -p 5432:5432 -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=mydb postgres:13
```
Create the `users` table:
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```
