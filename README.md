# 🌡️ go-sensor-data-collector

A lightweight Go microservice that exposes HTTP endpoints to receive and persist sensor readings (SHT31 and DS18B20) into TimescaleDB. Each sensor has its own dedicated POST route under `/sensor`. Timestamp (`ts`) is assigned automatically by the database with `DEFAULT now()`.

---

## ✨ Features

- 🔌 **Dedicated endpoints**

  - `POST /sensor/sht31` for SHT31 (temperature + humidity)
  - `POST /sensor/ds18b20` for DS18B20 (temperature only)
  - `GET /health` for health checks (no authentication required)

- 🏗️ **Repository Pattern Architecture**

  - Interface-based repository design for better testability
  - Dependency injection for clean separation of concerns
  - GORM with optimized connection pooling

- ⏱️ **Automatic timestamps**

  - Database column `ts TIMESTAMPTZ NOT NULL DEFAULT now()` ensures each row is stamped on insert.

- 🔐 **Token-based authentication (optional)**

  - If `API_TOKEN` is set in the environment, `/sensor/*` routes require `Authorization: Bearer <API_TOKEN>`
  - Health check endpoint remains public for monitoring

- 🐳 **Dockerized with Health Checks**
  - Multi-stage `Dockerfile` produces a minimal Alpine-based image
  - `docker-compose.yml` brings up TimescaleDB and the Go service with health monitoring
  - Environment variables managed via `.env` file

---

## 📋 Prerequisites

- Docker ≥ 20.10 or Podman ≥ 4.0
- Docker Compose ≥ 1.29 or podman-compose
- (Optional) `psql` client, if you want to query the database from your host
- (Optional) [Bruno](https://www.usebruno.com/) for API testing

---

## 🚀 Getting Started

1. **Clone the repository**

   ```bash
   git clone https://github.com/tu_usuario/go-sensor-data-collector.git
   cd go-sensor-data-collector
   ```

2. **Configure environment variables**

   ```bash
   cp .env.example .env
   # Edit .env with your preferred values
   ```

   Key variables:

   - `PORT` - Application port (default: 3000)
   - `API_TOKEN` - Optional Bearer token for authentication
   - `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB` - Database credentials
   - `DATABASE_URL` - Full PostgreSQL connection string

3. **Build and start services with Docker Compose**

   Using Docker:

   ```bash
   docker-compose up -d --build
   ```

   Using Podman:

   ```bash
   podman-compose up -d --build
   ```

   This will spin up:

   - **TimescaleDB** (`timescaledb:latest-pg14`), with an init script that creates two hypertables:
     - `sht31_readings (id, temperature, humidity, ts)`
     - `ds18b20_readings (id, temperature, ts)`
   - **Go microservice** (`go-sensor-data-collector`), listening on port 3000 (default).

   Both services include health checks:

   - TimescaleDB: `pg_isready` check every 5 seconds
   - App: `/health` endpoint check every 10 seconds with 40s startup grace period

4. **Verify both containers are running**

   Using Docker:

   ```bash
   docker-compose ps
   ```

   Using Podman:

   ```bash
   podman-compose ps
   ```

   You should see both services in the "Up (healthy)" state.

---

## ⚙️ Environment Variables

All environment variables are configured in the `.env` file:

### Application Variables

- `PORT` (default: `3000`)
  Port on which the Go HTTP server listens.

- `API_TOKEN` (optional)
  If set, the service requires `/sensor/*` endpoints to include:
  ```
  Authorization: Bearer <API_TOKEN>
  ```
  If `API_TOKEN` is empty or unset, sensor endpoints are public.
  Note: `/health` endpoint is always public.

### Database Variables

- `POSTGRES_USER` (default: `goapp`)
  PostgreSQL username for TimescaleDB.

- `POSTGRES_PASSWORD` (default: `secret`)
  PostgreSQL password for TimescaleDB.

- `POSTGRES_DB` (default: `sensors`)
  PostgreSQL database name.

- `POSTGRES_PORT` (default: `5432`)
  PostgreSQL port.

- `DATABASE_URL`
  Full PostgreSQL connection string. Example:
  ```
  host=timescaledb user=goapp password=secret dbname=sensors port=5432 sslmode=disable TimeZone=UTC
  ```

These variables are automatically injected into both `docker-compose.yml` services.

---

## 📡 API Endpoints

### 💚 Health Check

```
GET /health
```

**No authentication required** - This endpoint is public for health monitoring.

Returns:

```json
{
  "status": "healthy"
}
```

### POST /sensor/sht31

- **URL:** `/sensor/sht31`

- **Method:** `POST`

- **Headers:**

  - `Content-Type: application/json`
  - `Authorization: Bearer <API_TOKEN>` (if `API_TOKEN` is set)

- **Request Body:**

  ```json
  {
    "temperature": <float>,
    "humidity":    <float>
  }
  ```

  Both `temperature` and `humidity` are required.

- **Response (201 Created):**

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

  - The `datetime` field is the `ts` column assigned by the database.

- **Error Responses:**

  - `400 Bad Request` if JSON is invalid or missing required fields.
  - `401 Unauthorized` if `API_TOKEN` is set and the header is missing/invalid.
  - `500 Internal Server Error` on database errors.

### POST /sensor/ds18b20

- **URL:** `/sensor/ds18b20`

- **Method:** `POST`

- **Headers:**

  - `Content-Type: application/json`
  - `Authorization: Bearer <API_TOKEN>` (if `API_TOKEN` is set)

- **Request Body:**

  ```json
  {
    "temperature": <float>
  }
  ```

  `temperature` is required; `humidity` is omitted, since DS18B20 does not supply humidity.

- **Response (201 Created):**

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

- **Error Responses:**

  - `400 Bad Request` if JSON is invalid or missing `temperature`.
  - `401 Unauthorized` if `API_TOKEN` is set and the header is missing/invalid.
  - `500 Internal Server Error` on database errors.

---

## 🗄️ Database Schema (TimescaleDB)

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

- **`ts TIMESTAMPTZ NOT NULL DEFAULT now()`** ensures the database stamps each row with the server’s current time.
- Both tables are converted into TimescaleDB hypertables for efficient time-series storage and partitioning.

---

## 📁 Folder Structure

```
go-sensor-data-collector/
├── bruno/                          # Bruno API collection for testing
│   ├── bruno.json                  # Collection configuration
│   ├── environments/
│   │   └── Local.bru               # Local environment variables
│   ├── Health Check.bru            # GET /health request
│   ├── POST SHT31 Reading.bru      # POST /sensor/sht31 request
│   └── POST DS18B20 Reading.bru    # POST /sensor/ds18b20 request
├── cmd/
│   └── server/
│       └── main.go                 # Application entrypoint: DI, router, middleware
├── internal/
│   ├── config/
│   │   └── config.go               # Viper-based configuration loader
│   ├── db/
│   │   └── postgres.go             # DB initialization with connection pooling
│   ├── handlers/
│   │   ├── sht31.go                # SHT31 handler with repository injection
│   │   └── ds18b20.go              # DS18B20 handler with repository injection
│   ├── middleware/
│   │   └── auth.go                 # Optional Bearer token auth middleware
│   ├── models/
│   │   ├── sht31.go                # GORM model for sht31_readings
│   │   └── ds18b20.go              # GORM model for ds18b20_readings
│   └── repository/
│       ├── sensor_repository.go    # Repository interface definition
│       ├── sht31_repository.go     # SHT31 repository implementation
│       └── ds18b20_repository.go   # DS18B20 repository implementation
├── sql/
│   └── init-timescaledb.sql        # SQL script to create hypertables
├── .env.example                    # Example environment variables
├── .env                            # Your local environment variables (git-ignored)
├── go.mod                          # Module definition and dependencies
├── go.sum
├── Dockerfile                      # Multi-stage Docker build (Alpine-based)
└── docker-compose.yml              # Orchestrates TimescaleDB + Go service with health checks
```

---

## 🧪 Testing with Bruno

A Bruno API collection is included in the `bruno/` directory for easy testing:

1. **Install Bruno** from [usebruno.com](https://www.usebruno.com/)
2. **Open the collection** by pointing Bruno to the `bruno/` folder
3. **Configure environment** in `bruno/environments/Local.bru`:
   - Set `base_url` (default: `http://localhost:3000`)
   - Set `api_token` if using authentication
4. **Run requests** to test your endpoints

## 🏛️ Architecture

### Repository Pattern

The application uses an interface-based repository pattern:

- **Interfaces** (`internal/repository/sensor_repository.go`) define contracts
- **Implementations** (SHT31Repository, DS18B20Repository) handle data access
- **Handlers** receive repositories via dependency injection
- **Benefits**: Testability, separation of concerns, maintainability

### Database Connection Pooling

Optimized GORM configuration in `internal/db/postgres.go`:

- Max idle connections: 10
- Max open connections: 100
- Connection max lifetime: 1 hour
- Connection max idle time: 10 minutes
- Prepared statements enabled for better performance

## 🔧 Customization

- **Protecting Endpoints:**
  Set `API_TOKEN` in `.env` file. If unset, sensor endpoints are public. Health check is always public.

- **Changing Ports:**
  Modify `PORT` in `.env` file. The docker-compose.yml will automatically use the new port.

- **Using Another Database:**
  Replace `timescale/timescaledb:latest-pg14` with any PostgreSQL image. Remove hypertable logic if not using TimescaleDB.

---

## 🔍 Local Database Inspection

1. **Via `psql` on host:**

   ```bash
   psql "host=localhost port=5432 user=goapp password=secret dbname=sensors sslmode=disable"
   ```

2. **Via Docker Exec:**

   ```bash
   docker exec -it timescaledb psql -U goapp -d sensors
   ```

   Via Podman:

   ```bash
   podman exec -it timescaledb psql -U goapp -d sensors
   ```

3. **Example Queries:**

   ```sql
   -- View latest SHT31 readings
   SELECT id, temperature, humidity, ts
   FROM sht31_readings
   ORDER BY ts DESC
   LIMIT 5;

   -- View latest DS18B20 readings
   SELECT id, temperature, ts
   FROM ds18b20_readings
   ORDER BY ts DESC
   LIMIT 5;

   -- Count total readings
   SELECT
     (SELECT COUNT(*) FROM sht31_readings) as sht31_count,
     (SELECT COUNT(*) FROM ds18b20_readings) as ds18b20_count;
   ```

---

## 💻 Development

### Running Locally (without Docker)

1. Ensure you have PostgreSQL/TimescaleDB running
2. Copy `.env.example` to `.env` and configure `DATABASE_URL`
3. Run the application:
   ```bash
   go run ./cmd/server
   ```

### Building

```bash
go build -o sensor-collector ./cmd/server
./sensor-collector
```
