# 🛒 Microservice E-Commerce Platform

A modern, scalable e-commerce platform built with microservices architecture, featuring real-time notifications, payment processing, and comprehensive product management capabilities.

## 📋 Table of Contents

- [Architecture Overview](#architecture-overview)
- [Technology Stack](#technology-stack)
- [Services](#services)
- [Infrastructure](#infrastructure)
  - [Infrastructure as Code (Terraform)](#infrastructure-as-code-terraform)
  - [Helm Charts](#helm-charts)
- [Local Development (Docker Compose)](#local-development-docker-compose)
- [Production Deployment (Kubernetes)](#production-deployment-kubernetes)
- [Project Structure](#project-structure)

## 🏗️ Architecture Overview

This project follows a microservices architecture pattern with the following key components:

- **Backend Services**: Four independent Go-based microservices
- **API Gateway**: Traefik reverse proxy for routing and authentication
- **Frontend**: Vue.js 3 SPA with TypeScript
- **Infrastructure**: Kubernetes (EKS) deployment with Terraform
- **Communication**: RESTful APIs with Socket.IO for real-time features

### System Design

```
┌─────────────┐
│   Frontend  │ (Vue.js 3 + TypeScript)
│  (Port 5173)│
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────────┐
│        Traefik API Gateway              │
│  (Routing, CORS, Authentication)        │
│          (Port 81)                      │
└───┬────────┬────────┬────────┬─────────┘
    │        │        │        │
    ▼        ▼        ▼        ▼
┌────────┐┌─────────┐┌────────┐┌──────────────┐
│  User  ││ Product ││ Order  ││ Notification │
│Service ││ Service ││Service ││  Service     │
│:8185   ││  :8186  ││ :8187  ││   :8188      │
└────────┘└─────────┘└────────┘└──────────────┘
```

## 🛠️ Technology Stack

### Backend

- **Language**: Go 1.24+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Cloud Storage**: AWS S3 (for product images)
- **Payment Processing**: Stripe
- **Real-time Communication**: Socket.IO

### Frontend

- **Framework**: Vue.js 3.5+ with Composition API
- **Language**: TypeScript
- **UI Library**: Element Plus
- **Build Tool**: Vite
- **State Management**: Vue Router
- **Charts**: Chart.js with vue-chartjs
- **Maps**: Leaflet
- **HTTP Client**: Axios

### Infrastructure

- **Container Orchestration**: Kubernetes (Amazon EKS)
- **Infrastructure as Code**: Terraform
- **Reverse Proxy**: Traefik v3.6
- **Package Manager**: Helm
- **Container Runtime**: Docker

## 🚀 Services

### 1. User Service (Port 8185)

Handles user authentication, authorization, and profile management.

**Key Features:**
- User registration and login
- JWT-based authentication
- Role-based access control (RBAC)
- Profile management
- Forward authentication for other services

**API Endpoints:**
- Public: `/api/user/public/*` - Registration, login, token verification
- Private: `/api/user/*` - Profile management (requires authentication)

### 2. Product Service (Port 8186)

Manages product catalog, inventory, and search functionality.

**Key Features:**
- Product CRUD operations
- Category management
- Search and filtering
- Product recommendations based on search history
- Image upload to AWS S3
- Seller product management

**API Endpoints:**
- Public: `/api/product/public/*` - Product browsing, search
- Private: `/api/product/*` - Product management (seller/admin only)

### 3. Order Service (Port 8187)

Handles order processing, checkout, and payment integration.

**Key Features:**
- Shopping cart management
- Checkout session creation
- Stripe payment integration
- Order tracking and history
- Order status management

**API Endpoints:**
- Public: `/api/order/public/*` - Webhook handlers (Stripe)
- Private: `/api/order/*` - Order management, checkout

### 4. Notification Service (Port 8188)

Provides real-time notifications and messaging.

**Key Features:**
- Real-time Socket.IO connections
- Event-based notifications
- User notification management
- WebSocket support for live updates

**API Endpoints:**
- Socket.IO: `/socket.io/*` - Real-time connections
- Public: `/api/notification/public/*`
- Private: `/api/notification/*` - Notification management

## �️ Infrastructure

### Infrastructure as Code (Terraform)

The project uses Terraform to provision and manage AWS infrastructure. The IaC modules are located in the `iac/` directory.

#### 📦 Available Modules

**1. EKS Cluster (`iac/cluster/`)**

Provisions a complete Amazon EKS cluster with VPC, subnets, and node groups.

**Resources Created:**
- Custom VPC with public and private subnets across multiple availability zones
- EKS cluster with managed node groups (2-3 nodes, scalable)
- IAM roles and policies for cluster and worker nodes
- Load balancer controller permissions for NLB/ALB creation
- Security groups and network ACLs

**Configuration:**
```bash
cd iac/cluster
terraform init
terraform plan
terraform apply
```

**Key Variables:**
- `project`: Project name prefix
- `cidr_block`: VPC CIDR block
- `vpc_name`: VPC name
- `private_subnet_cidrs`: List of private subnet CIDRs
- `public_subnet_cidrs`: List of public subnet CIDRs
- `availability_zones`: AWS availability zones

**2. Frontend Hosting (`iac/frontend/`)**

Sets up S3 bucket and CloudFront distribution for Vue.js static site hosting.

**Resources Created:**
- S3 bucket with website configuration
- CloudFront distribution with Origin Access Identity (OAI)
- SSL/TLS certificate support (ACM)
- Optional Route 53 DNS records for custom domains
- Cache policies optimized for SPAs

**Features:**
- Automatic HTTPS redirect
- SPA routing support (404/403 → index.html)
- Gzip compression enabled
- Custom domain support with ACM certificates

**Configuration:**
```bash
cd iac/frontend
terraform init
terraform plan
terraform apply
```

**Key Variables:**
- `bucket_name`: S3 bucket name
- `project`: Project identifier
- `cloudfront_price_class`: CDN price class
- `custom_domain`: Optional custom domain
- `acm_certificate_arn`: SSL certificate ARN
- `route53_zone_id`: Route 53 hosted zone ID

**3. Traefik Load Balancer (`iac/traefik/`)**

Configures Network Load Balancer for Traefik ingress controller.

**Configuration:**
```bash
cd iac/traefik
terraform init
terraform apply
```

### Helm Charts

Each microservice has a dedicated Helm chart for Kubernetes deployment. Charts are located in the `helm/` directory.

#### 📋 Chart Structure

All service charts follow this structure:

```
helm/<service-name>/
├── Chart.yaml           # Chart metadata
├── values.yaml          # Default configuration values
├── values-secret.yaml   # Secret values (env vars, credentials)
└── templates/
    ├── deployment.yaml  # Kubernetes Deployment
    ├── service.yaml     # Kubernetes Service
    ├── configmap.yaml   # Configuration data
    ├── secret.yaml      # Sensitive data
    └── ingress.yaml     # Traefik IngressRoute
```

#### 🎯 Available Charts

**1. User Service (`helm/user-service/`)**
- Deployment with configurable replicas
- PostgreSQL database connection
- JWT secret management
- AWS credentials for S3 access
- Traefik routing with forward auth

**2. Product Service (`helm/product-service/`)**
- Product catalog deployment
- S3 integration for product images
- Database connection pooling
- Search and recommendation engine

**3. Order Service (`helm/order-service/`)**
- Order processing deployment
- Stripe API integration
- Payment webhook configuration
- Inter-service communication with User and Product services

**4. Notification Service (`helm/notification-service/`)**
- Real-time Socket.IO server
- WebSocket support
- Event-driven notifications
- Redis integration (optional)

**5. Traefik (`helm/traefik/`)**
- API Gateway deployment
- IngressRoute definitions
- CORS middleware configuration
- Forward authentication middleware
- Network Load Balancer service

#### ⚙️ Configuration

Each chart uses two values files:

**`values.yaml`** - Non-sensitive configuration:
```yaml
image:
  repository: your-docker-username/user-service
  tag: "v1.0.0"
  pullPolicy: IfNotPresent

replicaCount: 2

service:
  type: ClusterIP
  port: 8085
```

**`values-secret.yaml`** - Sensitive data (DO NOT commit):
```yaml
env:
  DB_HOST: "your-postgres-host"
  DB_PASSWORD: "your-db-password"
  JWT_SECRET: "your-jwt-secret"
  AWS_ACCESS_KEY_ID: "your-aws-key"
  STRIPE_SECRET_KEY: "your-stripe-key"
```

#### 🚀 Installation Commands

Install charts using the provided scripts:

```bash
# Install all services in order
./scripts/helm-install-traefik.sh
./scripts/helm-install-user-service.sh
./scripts/helm-install-product-service.sh
./scripts/helm-install-order-service.sh
./scripts/helm-install-notification-service.sh
```

Or manually:

```bash
helm install user-service helm/user-service \
  -f helm/user-service/values.yaml \
  -f helm/user-service/values-secret.yaml
```

#### 🔄 Updating Deployments

```bash
helm upgrade user-service helm/user-service \
  -f helm/user-service/values.yaml \
  -f helm/user-service/values-secret.yaml
```

#### 🗑️ Uninstalling

```bash
helm uninstall user-service
```

---

## 💻 Local Development (Docker Compose)

For local development, use Docker Compose to run all services on your machine.

### Prerequisites

- **Docker** and **Docker Compose**: For running containers
- **Node.js** v20.19+ or v22.12+: For frontend development
- **Go** 1.24+: For backend development (optional)

### Setup

#### 1. Clone the Repository

```bash
git clone <repository-url>
cd e-commerce
```

#### 2. Configure Environment Variables

Each service requires its own `.env` file. Create them in the respective service directories:

```bash
# services/user-service/.env
# services/product-service/.env
# services/order-service/.env
# services/notification-service/.env
# frontend/.env
```

Refer to `.env.example` files (if available) in each service directory for required variables.

#### 3. Start Services with Docker Compose

```bash
docker-compose up -d
```

This will start:
- Traefik API Gateway: http://localhost:81
- Traefik Dashboard: http://localhost:8081
- User Service: http://localhost:8185
- Product Service: http://localhost:8186
- Order Service: http://localhost:8187
- Notification Service: http://localhost:8188

#### 4. Start Frontend Development Server

```bash
cd frontend
npm install
npm run dev
```

Frontend will be available at: http://localhost:5173

### Accessing the Application

- **Frontend**: http://localhost:5173
- **API Gateway**: http://localhost:81
- **Traefik Dashboard**: http://localhost:8081/dashboard/

All API requests from the frontend should go through the Traefik gateway at port 81.

## 💻 Development

### Backend Development

Each Go service follows the same structure:

```
services/<service-name>/
├── cmd/           # Application entry point
├── internal/      # Private application code
│   ├── handler/   # HTTP handlers
│   ├── service/   # Business logic
│   ├── repository/# Data access layer
│   ├── model/     # Data models
│   └── middleware/# HTTP middlewares
├── Dockerfile
├── go.mod
└── .env
```

**Running a service locally:**

```bash
cd services/<service-name>
go mod download
go run cmd/main.go
```

**Running tests:**

```bash
go test ./...
```

### Frontend Development

**Directory Structure:**

```
frontend/
├── src/
│   ├── components/  # Reusable Vue components
│   ├── views/       # Page components
│   ├── router/      # Vue Router configuration
│   ├── api/         # API client modules
│   ├── types/       # TypeScript type definitions
│   └── assets/      # Static assets
├── public/          # Public static files
└── index.html       # Entry HTML
```

**Available Scripts:**

```bash
npm run dev        # Start development server
npm run build      # Build for production
npm run preview    # Preview production build
npm run type-check # Run TypeScript type checking
npm run format     # Format code with Prettier
```

### Building Docker Images

Use the provided build script to build and push Docker images:

```bash
./scripts/build.sh <DOCKER_USERNAME> <VERSION> <SERVICE_NAME>
```

Example:

```bash
./scripts/build.sh myusername v1.0.0 user-service
```

### CORS Configuration

The application is configured to allow CORS from the following origins:

- `http://localhost:3000`
- `http://localhost:8080`
- `http://localhost:5173` (Frontend dev server)
- `http://localhost:24000`
- `http://djj0amhpbxjv0.cloudfront.net` (CloudFront distribution)

Modify CORS settings in `docker-compose.yaml` under the Traefik labels:

```yaml
- "traefik.http.middlewares.cors-headers.headers.accesscontrolalloworiginlist=http://localhost:5173,..."
```

### Troubleshooting

**Services not starting:**
```bash
docker-compose logs <service-name>
```

**Reset everything:**
```bash
docker-compose down -v  # Removes containers and volumes
docker-compose up -d
```

**Rebuild specific service:**
```bash
docker-compose up -d --build <service-name>
```

---

## 🚢 Production Deployment (Kubernetes)

Deploy the application to a Kubernetes cluster using Helm charts and Terraform.

### Prerequisites

- **AWS Account**: For EKS cluster and related resources
- **kubectl**: Kubernetes CLI tool
- **Helm**: v3+ for package management
- **Terraform**: v1.0+ for infrastructure provisioning
- **AWS CLI**: Configured with appropriate credentials

### Deployment Steps

#### 1. Provision Infrastructure

**Step 1:** Create VPC and EKS Cluster

```bash
cd iac/cluster
terraform init
terraform plan
terraform apply
```

This creates:
- VPC with public/private subnets
- EKS cluster with 2-3 worker nodes
- IAM roles and policies
- Security groups

**Step 2:** Configure kubectl

```bash
aws eks update-kubeconfig --region <region> --name <cluster-name>
kubectl get nodes  # Verify connection
```

**Step 3:** Deploy Frontend to S3/CloudFront

```bash
cd iac/frontend
terraform init
terraform apply

# Build and upload frontend
cd ../../frontend
npm run build
aws s3 sync dist/ s3://<bucket-name>/
```

#### 2. Build and Push Docker Images

Build all service images and push to Docker Hub (or your registry):

```bash
# Build all services
./scripts/build.sh <DOCKER_USERNAME> <VERSION> user-service
./scripts/build.sh <DOCKER_USERNAME> <VERSION> product-service
./scripts/build.sh <DOCKER_USERNAME> <VERSION> order-service
./scripts/build.sh <DOCKER_USERNAME> <VERSION> notification-service
```

Example:
```bash
./scripts/build.sh myusername v1.0.0 user-service
```

#### 3. Configure Helm Values

Update `values-secret.yaml` for each service with your credentials:

```yaml
# helm/user-service/values-secret.yaml
env:
  DB_HOST: "prod-postgres.xxxxx.rds.amazonaws.com"
  DB_USER: "admin"
  DB_PASSWORD: "your-secure-password"
  JWT_SECRET: "your-jwt-secret"
  AWS_ACCESS_KEY_ID: "AKIA..."
  AWS_SECRET_ACCESS_KEY: "..."
```

> ⚠️ **Security:** Never commit `values-secret.yaml` files. Add them to `.gitignore`.

#### 4. Deploy Services with Helm

Deploy services in order:

```bash
# 1. Deploy Traefik (API Gateway)
./scripts/helm-install-traefik.sh

# 2. Deploy backend services
./scripts/helm-install-user-service.sh
./scripts/helm-install-product-service.sh
./scripts/helm-install-order-service.sh
./scripts/helm-install-notification-service.sh
```

#### 5. Verify Deployment

```bash
# Check all pods are running
kubectl get pods

# Check services
kubectl get svc

# Check Traefik routes
kubectl get ingressroute

# Get Load Balancer URL
kubectl get svc traefik -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
```

#### 6. Configure DNS (Optional)

Point your domain to the Load Balancer:

```bash
# Get LB hostname
LB_HOST=$(kubectl get svc traefik -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')

# Create Route 53 record (or use your DNS provider)
aws route53 change-resource-record-sets --hosted-zone-id <ZONE_ID> ...
```

### Updating Deployments

**Update a specific service:**

```bash
# 1. Build new image
./scripts/build.sh <username> <new-version> user-service

# 2. Update values.yaml with new image tag
vim helm/user-service/values.yaml  # Update image.tag

# 3. Upgrade Helm release
helm upgrade user-service helm/user-service \
  -f helm/user-service/values.yaml \
  -f helm/user-service/values-secret.yaml
```

**Rollback a deployment:**

```bash
helm rollback user-service
```

### Monitoring and Logs

**View pod logs:**
```bash
kubectl logs -f <pod-name>
```

**View all service logs:**
```bash
kubectl logs -l app=user-service
```

**Describe pod issues:**
```bash
kubectl describe pod <pod-name>
```

### Scaling

**Scale a deployment:**
```bash
kubectl scale deployment user-service --replicas=5
```

**Or update Helm values:**
```yaml
# values.yaml
replicaCount: 5
```

```bash
helm upgrade user-service helm/user-service -f helm/user-service/values.yaml
```

## 📁 Project Structure

```
e-commerce/
├── services/                    # Backend microservices
│   ├── user-service/           # User authentication & management
│   ├── product-service/        # Product catalog & search
│   ├── order-service/          # Order processing & payments
│   └── notification-service/   # Real-time notifications
├── frontend/                    # Vue.js frontend application
├── helm/                        # Kubernetes Helm charts
│   ├── user-service/
│   ├── product-service/
│   ├── order-service/
│   ├── notification-service/
│   └── traefik/
├── iac/                         # Infrastructure as Code (Terraform)
│   ├── cluster/                # EKS cluster provisioning
│   ├── frontend/               # Frontend hosting (S3/CloudFront)
│   └── traefik/                # Load balancer configuration
├── scripts/                     # Deployment and build scripts
│   ├── build.sh                # Docker build script
│   └── helm-install-*.sh       # Helm deployment scripts
├── data-process/                # Data processing utilities
├── docker-compose.yaml          # Local development environment
├── TRAEFIK-LOCAL-SETUP.md      # Traefik configuration guide
└── docker-compose-cors-guide.md # CORS configuration guide
```

