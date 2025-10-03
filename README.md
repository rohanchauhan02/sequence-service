
# Sequence Service

A simple Golang API for managing email sequences.
Each sequence can have multiple steps (email subject + content), with options for open and click tracking.

---

## Features

* Create a sequence with steps
* Update a sequence step
* Delete a sequence step
* Update sequence tracking settings
* Data stored in PostgreSQL

---

## Setup

### Requirements

* Go 1.23+
* PostgreSQL

### Run locally

1. Clone the repo
2. Set up `configs/app.config.local.yml` file with your config settings:
3. Install dependencies:

   ```bash
   go mod tidy
   ```

4. Start the server:

   ```bash
   go run main.go
   ```

Server will run on: **[http://localhost:8080](http://localhost:8080)**

---

## Example API

**Check Service Health**
`GET /api/v1/health`

---
