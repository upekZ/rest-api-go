# REST API for CRUD operations using Go

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Overview
A scalable REST API for managing users, built with Go. It supports CRUD operations, email and mobile number uniqueness validation. The API uses PostgreSQL for persistent storage.

### Key Features
- **User CRUD**: Create, read, update, and delete users with validated email and mobile number uniqueness.
- **Tech Stack**: Go, Chi, Pgx (PostgreSQL).

## Prerequisites
- Go 1.20+
- PostgreSQL 13+
- Docker (for running services)

## Setup

### 1. Clone the Repository
```bash
git clone https://github.com/upekZ/rest-api-go.git
cd rest-api-go
```

### 2. Dependencies
github.com/go-chi/chi/v5
github.com/jackc/pgx/v5
github.com/gorilla/websocket
github.com/patrickmn/go-cache

### 3. Configure Environment
Create a `.env` file:
```env
DATABASE_URL=postgres://user:password@localhost:5432/mydb
```
Start PostgreSQL instance in Docker with docker-compose.yml:
```env
docker-compose up -d
```

Create the `users` table:
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_status AS ENUM ('Active', 'Inactive');

CREATE TABLE IF NOT EXISTS "user" (
                                      userId uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name varchar(100) NOT NULL,
    last_name varchar(100) NOT NULL,
    email varchar(100) NOT NULL,
    phone varchar(100),
    age integer,
    "status" user_status DEFAULT 'Active'
    );
```


### ToDos

Scripts to create user table in postgre instace
complete Unit and integration tests
