# go-sensor-data-collector

A lightweight Go microservice that exposes HTTP endpoints to receive and persist sensor readings (SHT31 and DS18B20) into TimescaleDB. Each sensor has its own dedicated POST route under `/sensor`. Timestamp (`ts`) is assigned automatically by the database with `DEFAULT now()`.

---

## Features

* **Dedicated endpoints**

  * `POST /sensor/sht31` for SHT31 (temperature + humidity)
  * `POST /sensor/ds18b20` for DS18B20 (temperature only)

* **Automatic timestamps**

  * Database column `ts TIMESTAMPTZ NOT NULL DEFAULT now()` ensures each row is stamped on insert.

* **Token-based authentication (optional)**

  * If `API_TOKEN` is set in the environment, all requests must include header `Authorization: Bearer <API_TOKEN>`.

* **Dockerized**

  * Multi-stage `Dockerfile` produces a minimal static binary (`scratch` or Distroless).
  * `docker-compose.yml` brings up TimescaleDB and the Go service together.

---

## Prerequisites

* Docker ≥ 20.10
* Docker Compose ≥ 1.29
* (Optional) `psql` client, if you want to query the database from your host

---

## Getting Started

1. **Clone the repository**

   ```bash
   git clone https://github.com/tu_usuario/go-sensor-data-collector.git
   cd go-sensor-data-collector
   ```

2. **Build and start services with Docker Compose**

   ```bash
   docker-compose up -d --build
   ```

   This will spin up:

   * **TimescaleDB** (`timescaledb:latest-pg14`), with an init script that creates two hypertables:

     * `sht31_readings (id, temperature, humidity, ts)`
     * `ds18b20_readings (id, temperature, ts)`
   * **Go microservice** (`go-sensor-data-collector`), listening on port 5000.

3. **Verify both containers are running**

   ```bash
   docker-compose ps
   ```

   You should see `timescaledb` and `go-sensor-data-collector` in the “Up” state.

---

## Environment Variables

The Go service reads the following variables at runtime:

* `DATABASE_URL`
  Full PostgreSQL connection string. Defaults to:

  ```
  host=timescaledb user=goapp password=secret dbname=sensors port=5432 sslmode=disable TimeZone=UTC
  ```

  (Matches the `docker-compose.yml` service `timescaledb` settings.)

* `API_TOKEN` (optional)
  If set, the service requires every request to include:

  ```
  Authorization: Bearer <API_TOKEN>
  ```

  If `API_TOKEN` is empty or unset, all endpoints are public.

* `PORT` (optional)
  Port on which the Go HTTP server listens. Defaults to `5000`.

---

## API Endpoints

### Health-check

```
GET /sensor/
```

Returns:

```json
{ "message": "go-sensor-data-collector is up" }
```

### POST /sensor/sht31

* **URL:** `/sensor/sht31`

* **Method:** `POST`

* **Headers:**

  * `Content-Type: application/json`
  * `Authorization: Bearer <API_TOKEN>` (if `API_TOKEN` is set)

* **Request Body:**

  ```json
  {
    "temperature": <float>,
    "humidity":    <float>
  }
  ```

  Both `temperature` and `humidity` are required.

* **Response (201 Created):**

  ```json
  {
    "message": "SHT31 reading saved",
    "record": {
      "id":          <int>,
      "temperature": <float>,
      "humidity":    <float>,
      "datetime":    "<RFC3339 timestamp>"
    }
  }
  ```

  * The `datetime` field is the `ts` column assigned by the database.

* **Error Responses:**

  * `400 Bad Request` if JSON is invalid or missing required fields.
  * `401 Unauthorized` if `API_TOKEN` is set and the header is missing/invalid.
  * `500 Internal Server Error` on database errors.

### POST /sensor/ds18b20

* **URL:** `/sensor/ds18b20`

* **Method:** `POST`

* **Headers:**

  * `Content-Type: application/json`
  * `Authorization: Bearer <API_TOKEN>` (if `API_TOKEN` is set)

* **Request Body:**

  ```json
  {
    "temperature": <float>
  }
  ```

  `temperature` is required; `humidity` is omitted, since DS18B20 does not supply humidity.

* **Response (201 Created):**

  ```json
  {
    "message": "DS18B20 reading saved",
    "record": {
      "id":          <int>,
      "temperature": <float>,
      "datetime":    "<RFC3339 timestamp>"
    }
  }
  ```

* **Error Responses:**

  * `400 Bad Request` if JSON is invalid or missing `temperature`.
  * `401 Unauthorized` if `API_TOKEN` is set and the header is missing/invalid.
  * `500 Internal Server Error` on database errors.

---

## Database Schema (TimescaleDB)

The init script at `sql/init-timescaledb.sql` runs automatically on container startup:

```sql
-- Create SHT31 hypertable
CREATE TABLE IF NOT EXISTS sht31_readings (
    id          SERIAL       PRIMARY KEY,
    temperature DOUBLE PRECISION NOT NULL,
    humidity    DOUBLE PRECISION NOT NULL,
    ts          TIMESTAMPTZ  NOT NULL DEFAULT now()
);
SELECT create_hypertable('sht31_readings', 'ts', if_not_exists => TRUE);

-- Create DS18B20 hypertable
CREATE TABLE IF NOT EXISTS ds18b20_readings (
    id          SERIAL       PRIMARY KEY,
    temperature DOUBLE PRECISION NOT NULL,
    ts          TIMESTAMPTZ  NOT NULL DEFAULT now()
);
SELECT create_hypertable('ds18b20_readings', 'ts', if_not_exists => TRUE);
```

* **`ts TIMESTAMPTZ NOT NULL DEFAULT now()`** ensures the database stamps each row with the server’s current time.
* Both tables are converted into TimescaleDB hypertables for efficient time-series storage and partitioning.

---

## Folder Structure

```
go-sensor-data-collector/
├── cmd/
│   └── server/
│       └── main.go          # Application entrypoint: sets up router, middleware, and handlers
├── internal/
│   ├── db/
│   │   └── postgres.go      # DB initialization (GORM + Postgres)
│   ├── handlers/
│   │   ├── sht31.go         # Handler for POST /sensor/sht31
│   │   └── ds18b20.go       # Handler for POST /sensor/ds18b20
│   ├── middleware/
│   │   └── auth.go          # Optional token auth middleware
│   └── models/
│       ├── sht31.go         # GORM model for sht31_readings
│       └── ds18b20.go       # GORM model for ds18b20_readings
├── go.mod                   # Module definition and dependencies
├── go.sum
├── Dockerfile               # Multi-stage Docker build for Go binary
├── docker-compose.yml       # Orchestrates TimescaleDB + Go service
└── sql/
    └── init-timescaledb.sql # SQL script to create hypertables for SHT31 & DS18B20
```

---

## Customization

* **Protecting Endpoints:**
  Set `API_TOKEN` in `docker-compose.yml` (or your environment). If unset, endpoints are public.

* **Changing Ports:**

  * Modify `PORT` in `docker-compose.yml` or set it via environment.
  * Adjust the `EXPOSE` instruction in `Dockerfile` if you use a different port.

* **Using Another Database:**
  Replace `timescale/timescaledb:latest-pg14` with any PostgreSQL image. Remove hypertable logic if not using TimescaleDB.

---

## Local Database Inspection

1. **Via `psql` on host:**

   ```bash
   psql "host=localhost port=5432 user=goapp password=secret dbname=sensors sslmode=disable"
   ```
2. **Via Docker Exec:**

   ```bash
   docker exec -it timescaledb psql -U goapp -d sensors
   ```
3. **Example Query:**

   ```sql
   SELECT id, temperature, humidity, ts
   FROM sht31_readings
   ORDER BY ts DESC
   LIMIT 5;
   ```

---

