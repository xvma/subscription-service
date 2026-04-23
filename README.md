# Subscription Service

REST service for aggregating user online subscription data.

## Description

The service provides an HTTP API for managing user subscriptions with the following capabilities:
- CRUD operations on subscription records
- Calculation of total subscription cost for a selected period with filtering

## Tech Stack

- *Go 1.21* - Main programming language
- *Gin* - Web framework
- *PostgreSQL 15* - Database
- *Docker & Docker Compose* - Containerization and orchestration
- *Logrus* - Structured logging
- *Swagger* - API documentation

## Installation and Setup

### Prerequisites

- Docker and Docker Compose

### Run with Docker Composf

```bash
git clone <repository-url>
cd subscription-service
docker-compose up -d
docker-compose ps
docker-compose logs -f app
```

### Local Setup (without Docker)

```bash
go mod download
createdb subscriptions
psql -d subscriptions -f migrations/001_create_subscriptions_table.up.sql
go run cmd/main.go
```

## Configuration

Environment variables (.env file):

```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions
DB_SSLMODE=disable
SERVER_PORT=8080
LOG_LEVEL=info
```

## API Documentation

Base URL: `http://localhost:8080/api/v1`

### 1. Create Subscription

```http
POST /subscriptions
```

Request Body:
```json
{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2025"
}
```

Response (201 Created):
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025",
  "end_date": "12-2025",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### 2. Get Subscription by ID

```http
GET /subscriptions/{id}
```

### 3. Update Subscription

```http
PUT /subscriptions/{id}
Content-Type: application/json
```

Request Body:

```json
{
  "price": 500,
  "end_date": "06-2026"
}
```

### 4. Delete Subscription

```http
DELETE /subscriptions/{id}
```

### 5. Calculate Total Cost

```http
POST /subscriptions/total-cost
Content-Type: application/json
```

Request Body:
```json
{
  "start_date": "01-2024",
  "end_date": "12-2024",
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "service_name": "Yandex Plus"
}
```

Response:

```json
{
  "total_cost": 4800
}
```

### 6. Health Check

```http
GET /health
```
Response:
```json
{
  "status": "healthy"
}
```

## Database Schema

| Field | Type | Description |
|--------|---------|--------------|
| id | UUIT | Unique identifier |
| service_name | VARCHAR(255) | Service name |
| price | INTEGER | Cost in rubles |
| user_id | UUIT | User ID |
| start_date | DATE | Start date (MM-YYYY) |
| end_date | DATE | End date (optional) |
| created_at | TIMESTAMP | Creation date |
| updated_at | TIMESTAMP | Last update |

## Testing with cURL

```bash
# Create subscription
crl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Netflix",
    "price": 799,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "01-2024"
  }'

# Get subscription
crl http://localhost:8080/api/v1/subscriptions/123e4567-e89b-12d3-a456-426614174000

# Update subscription
crl -X PUT http://localhost:8080/api/v1/subscriptions/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{"price": 899}'

# Delete subscription
crl -X DELETE http://localhost:8080/api/v1/subscriptions/123e4567-e89b-12d3-a456-426614174000

# Calculate total cost
crl -X POST http://localhost:8080/api/v1/subscriptions/total-cost \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "01-2024",
    "end_date": "12-2024"
  }'
```

## Troubleshooting

```bash
# Check logs
docker-compose logs app

# Check database connection
docker-compose exec postgres pg_isready - postgres

# Apply migrations manually
docker-compose exec postgres psql -U postgres -d subscriptions \
  -f /docker-entrypoint-initdb/001_create_subscriptions_table.up.sql
```

## License

MIT License
