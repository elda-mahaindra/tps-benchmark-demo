# py-core - Balance Inquiry Service (Python 2.7)

A TCP-based service for account balance inquiry, intentionally built with deprecated Python 2.7 to demonstrate performance differences with modern Go implementations.

**Note**: Now uses TCP protocol with persistent connections instead of REST API.

## ⚠️ WARNING

This service uses **Python 2.7** which reached End of Life on January 1, 2020. This is intentionally used for demonstration purposes only to show the performance gap with modern technology stacks.

## Features

- **TCP Server**: Persistent connection protocol with generic request/response format
- **Message ID Correlation**: Async request/response handling
- **No Connection Pooling**: Creates new database connection per request (intentionally inefficient)
- **Python 2.7**: Uses deprecated Python version
- **Generic Operations**: Supports multiple operations via single protocol

## TCP Protocol

### Connection

- **Host**: 0.0.0.0
- **Port**: 5001
- **Type**: Persistent connection (stays open)

### Message Format

**Wire Protocol**:

```
[2-byte length (big-endian)][JSON payload]
```

### Request Format

```json
{
  "id_message": "unique-message-id",
  "operation": "operation_name",
  "params": {
    // operation-specific parameters
  }
}
```

### Response Format

```json
{
  "id_message": "unique-message-id",
  "status": "000",
  "err_info": "",
  "data": {
    // operation-specific data
  }
}
```

**Status Codes**:

- `000`: Success
- `999`: Error
- `998`: Timeout
- `997`: Request dropped

## Supported Operations

### 1. Ping

**Request**:

```json
{
  "id_message": "msg-123",
  "operation": "ping",
  "params": {}
}
```

**Response**:

```json
{
  "id_message": "msg-123",
  "status": "000",
  "err_info": "",
  "data": {
    "message": "pong",
    "service": "py-core",
    "timestamp": "2025-10-20T12:00:00Z"
  }
}
```

### 2. Get Account by Account Number

**Request**:

```json
{
  "id_message": "msg-456",
  "operation": "get_account_by_account_number",
  "params": {
    "account_number": "1001000000001"
  }
}
```

**Response**:

```json
{
  "id_message": "msg-456",
  "status": "000",
  "err_info": "",
  "data": {
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
}
```

## Running with Docker

### Build and Run

```bash
# From project root
docker-compose -f docker-compose.dev.yml up py-core

# Or rebuild
docker-compose -f docker-compose.dev.yml build py-core
docker-compose -f docker-compose.dev.yml up py-core
```

### Test the Service

**Using the test client**:

```bash
# Inside the container
docker exec -it py-core python test_tcp_client.py localhost 5001

# Or from host (if port is exposed)
cd py-core
python test_tcp_client.py localhost 5001
```

**Using netcat (manual test)**:

```bash
# Connect to server
nc localhost 5001

# Send request (you need to calculate length manually)
# This is complex - use test_tcp_client.py instead
```

## File Structure

```
py-core/
├── app_tcp.py             # TCP server entry point
├── tcp_server.py          # TCP server implementation
├── request_handler.py     # Request routing and handling
├── database.py            # Database operations
├── test_tcp_client.py     # Test client for verification
├── config.json            # Configuration
├── requirements.txt       # Python dependencies
├── Dockerfile             # Docker build file
└── README.md              # This file
```

## Performance Characteristics

This service is intentionally slow due to:

1. **Python 2.7**: Old interpreter with no modern optimizations
2. **No Connection Pooling**: Creates/destroys DB connection per request
3. **Synchronous I/O**: Blocking operations (though threading helps)
4. **TCP Overhead**: Manual protocol implementation
5. **Thread-based Concurrency**: Less efficient than Go's goroutines

These characteristics make it perfect for demonstrating the performance benefits of:

- Modern language runtimes (Go 1.25+)
- Connection pooling
- Efficient concurrency models (goroutines vs threads)
- Optimized database drivers

## Configuration

Configuration is stored in `config.json`:

```json
{
  "tcp_host": "0.0.0.0",
  "tcp_port": 5001,
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

- psycopg2-binary 2.8.6
- Werkzeug 1.0.1

## Development

### Run Locally (Not Recommended)

```bash
# Install dependencies (Python 2.7 required)
pip install -r requirements.txt

# Run TCP server
python app_tcp.py

# In another terminal, test
python test_tcp_client.py localhost 5001
```

### Debug Mode

```python
# In app_tcp.py, set logging level to DEBUG
logging.basicConfig(level=logging.DEBUG)
```

## Architecture

### Persistent Connection Pattern

```
Client                          py-core Server
  |                                   |
  |---- TCP Connect ----------------->|
  |                                   |
  |===== Connection Maintained =======|
  |                                   |
  |---- Request 1 (with id_msg_1) --->|
  |                                   |--- Process in thread
  |<--- Response 1 (id_msg_1) --------|
  |                                   |
  |---- Request 2 (with id_msg_2) --->|
  |                                   |--- Process in thread
  |<--- Response 2 (id_msg_2) --------|
  |                                   |
  |===== Connection Maintained =======|
  |                                   |
  |---- Disconnect ------------------->|
```

### Key Differences from REST

| Aspect              | REST (Old)         | TCP (New)                |
| ------------------- | ------------------ | ------------------------ |
| Connection          | New per request    | Persistent               |
| Request Correlation | Implicit           | Via message ID           |
| Concurrency         | Serial             | Async via threading      |
| Protocol            | HTTP/JSON          | TCP/Length-prefixed JSON |
| Operations          | Multiple endpoints | Generic operation field  |

## Comparison with Go Implementation

While both handle the same business logic, the implementations differ:

**Go (go-core)**:

- gRPC with specific RPCs
- Connection pooling
- Goroutine-based concurrency
- Type-safe Protocol Buffers

**Python (py-core)**:

- TCP with generic operations
- No connection pooling
- Thread-based concurrency
- JSON serialization

Both query the same PostgreSQL database and return identical data structures.
