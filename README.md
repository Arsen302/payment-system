# Payment Processing System

A scalable, microservice-based payment processing system built with Go, gRPC, Kafka, PostgreSQL, and Redis.

## Architecture

The system consists of three microservices:

1. **Auth Service**: Handles user authentication and authorization using JWT tokens.
2. **Payment Service**: Processes payments and publishes events to Kafka.
3. **Notification Service**: Consumes payment events from Kafka and sends notifications.

## Technologies

- **Go 1.24**: Modern, concurrent programming language
- **gRPC**: High-performance RPC framework
- **Kafka**: Event streaming platform for high-throughput messaging
- **PostgreSQL**: Relational database with sharding capabilities
- **Redis**: In-memory data store for caching and session management
- **Docker & Docker Compose**: Containerization and orchestration

## Getting Started

### Prerequisites

- Go 1.24+
- Docker and Docker Compose
- Make (optional)

### Environment Setup

The project uses environment variables for configuration. Each service has its own environment files:

1. Copy the example environment files to create your local configuration:
   ```bash
   cp .env.example .env
   cp auth-service/.env.example auth-service/.env
   cp payment-service/.env.example payment-service/.env
   cp notification-service/.env.example notification-service/.env
   ```

2. Edit the `.env` files to set your local configuration values.

3. The system supports different environment configurations:
   - `.env`: Default environment file
   - `.env.local`: Local development environment
   - `.env.production`: Production environment

### Running the Services

1. Clone the repository:
   ```bash
   git clone https://github.com/Arsen302/payment-system.git
   cd payment-system
   ```

2. Start the infrastructure and services using Docker Compose:
   ```bash
   docker-compose up
   ```

3. The services will be available at:
   - Auth Service: `localhost:50051`
   - Payment Service: `localhost:50052`

4. For running individual services without Docker:
   ```bash
   # Using the default .env file
   make run-auth

   # Using a specific environment file
   make run-auth ENV_FILE=.env.local
   ```

## Development

### Generating gRPC Code

To regenerate the gRPC code after modifying proto files:

```bash
make proto
```

### Building Individual Services

```bash
# Build all services
make build

# Or build specific services
cd auth-service && go build -o bin/auth-service ./cmd/server
```

## Project Structure

```
payment-system/
├── auth-service/          # Authentication and authorization service
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   └── api/
├── payment-service/       # Payment processing service
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   └── api/
├── notification-service/  # Notification handling service
│   ├── cmd/
│   ├── internal/
│   └── pkg/
├── proto/                 # Proto definitions
│   ├── auth.proto
│   └── payment.proto
├── deployments/           # Deployment configurations
│   └── postgres/
├── .env.example           # Example environment variables
├── .env.local             # Local development environment variables
├── .env.production        # Production environment variables
└── docker-compose.yml     # Docker Compose configuration
```