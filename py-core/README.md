# py-core - Balance Inquiry Service (Python 2.7)

A simple REST API service for account balance inquiry, intentionally built with deprecated Python 2.7 to demonstrate performance differences with modern implementations.

## ⚠️ WARNING
This service uses **Python 2.7** which reached End of Life on January 1, 2020. This is intentionally used for demonstration purposes only to show the performance gap with modern technology stacks.

## Features

- **REST API**: Simple HTTP/JSON interface
- **No Connection Pooling**: Creates new database connection per request (intentionally inefficient)
- **Python 2.7**: Uses deprecated Python version
- **Single File**: Simple, straightforward implementation

## API Endpoints

### Health Check
```
GET /health
```

Response:
```json
{
  "status": "ok",
  "service": "py-core"
}
```

### Get Account by Account Number
```
GET /api/v1/accounts/{account_number}
```

Example:
```bash
curl http://localhost:5000/api/v1/accounts/1001000000001
```

Response:
```json
{
  "account": {
    "account_id": 1,
    "account_number": "1001000000001",
    "customer_id": 1,
    "account_type": "WADIAH",
    "account_status": "ACTIVE",
    "balance": "15750000.00",
    "currency": "IDR",
    "opened_date": "2025-01-01",
    "closed_date": "",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  },
  "customer": {
    "customer_number": "CUST0000001",
    "full_name": "Ahmad Hidayat",
    "id_number": "3201012345670001",
    "phone_number": "081234567801",
    "email": "ahmad.hidayat@email.com",
    "address": "Jl. Merdeka No. 123, Jakarta",
    "date_of_birth": "1985-03-15"
  }
}
```

## Running with Docker

### Build and Run
```bash
# From project root
docker-compose -f docker-compose.dev.yml up py-core
```

### Test the Service
```bash
# Health check
curl http://localhost:5000/health

# Get account
curl http://localhost:5000/api/v1/accounts/1001000000001
```

## Performance Characteristics

This service is intentionally slow due to:

1. **Python 2.7**: Old interpreter with no modern optimizations
2. **No Connection Pooling**: Creates/destroys DB connection per request
3. **Synchronous I/O**: Blocking operations throughout
4. **Flask Overhead**: Additional framework overhead in Python 2.7

These characteristics make it perfect for demonstrating the performance benefits of:
- Modern language runtimes (Go 1.25+)
- Connection pooling
- Efficient concurrency models
- Optimized database drivers

## Configuration

Configuration is stored in `config.json`:

```json
{
  "host": "0.0.0.0",
  "port": 5000,
  "database": {
    "host": "postgres",
    "port": 5432,
    "database": "demo_db",
    "user": "postgres",
    "password": "changeme"
  }
}
```

## Dependencies

See `requirements.txt` for Python 2.7 compatible packages:
- Flask 1.1.4
- psycopg2-binary 2.8.6
- Werkzeug 1.0.1
