# REST API for CRUD operations using Go

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8)](https://golang.org/)

## Overview
A REST API for managing users, built with Go. It supports CRUD operations, email and mobile number uniqueness validation. The API uses PostgreSQL for persistent storage.

### Key Features
- **User CRUD**: Create, read, update, and delete users with validated email and mobile number uniqueness.
- **Tech Stack**: Go, Chi, Pgx (PostgreSQL), Gorilla.

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
- github.com/go-chi/chi/v5
- github.com/jackc/pgx/v5
- github.com/gorilla/websocket
- github.com/patrickmn/go-cache

### 3. Launch API
- Direct to rest-api-go
- 



### ToDos

- WebSocket Testing and instructions on running manual test
- Implement websocket broadcast to User deletion and updating
- Scripts to create user table in postgre instance
- complete Unit and integration tests
- extend caching for user retrieval
- Websocket accepts any connection without validation. Could be improved with validation check