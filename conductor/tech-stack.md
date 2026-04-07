# Tech Stack

## Programming Language
- **Go 1.25+:** Primary language for backend services.
  - *Why:* Performance, concurrency model, and strong typing.

## Frameworks & Libraries
- **Gin-Gonic:** High-performance HTTP framework.
- **GORM:** Object-Relational Mapper for PostgreSQL interaction.
- **API Gateway Authorizer:** Delegated authentication with identity forwarding via headers (X-User-Id, X-User-Email, X-User-Role).
- **Go-Playground/Validator:** Input validation for API requests.

## Database & Persistence
- **PostgreSQL 15:** Robust relational database for transactional data.
- **Docker Compose:** Local orchestration for the database and application.

## Infrastructure & DevOps
- **Docker:** Containerization of the API and its environment.
- **Kubernetes (Kind):** Local cluster simulation for production parity.
- **AWS (EKS):** Target platform for scalable deployments, integrated with GitHub Actions.
- **Secret Management:** Integration with GitHub Secrets for production and `.env` fallback for local development.

## Architecture
- **Hexagonal Architecture:** Decoupling business logic from infrastructure using Ports and Adapters.
