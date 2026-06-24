# Deployment Guide — AWS (EKS + CloudFront)

> Last updated: Jun 23, 2026
> Region: `ap-southeast-1`

---

## Architecture Overview

```
Internet
  │
  ├── CloudFront (d11unvcwv8df4t.cloudfront.net)
  │     ├── /api/*        → NLB (HTTP:80) → Traefik → K8s services
  │     ├── /socket.io/*  → NLB (HTTP:80) → Traefik → notification-service
  │     └── /*            → S3 (e-commerce-frontend-web) — static Vue.js
  │
  └── EKS Cluster (e-commerce-eks, k8s 1.33)
        VPC: 10.0.0.0/16
        Public subnets:  10.0.2.0/24 (1a), 10.0.4.0/24 (1b)  ← Traefik NLB
        Private subnets: 10.0.1.0/24 (1a), 10.0.3.0/24 (1b)  ← EKS nodes (t3.medium × 2)
```

**External dependencies (not in cluster):**
- MongoDB Atlas: `e-commerce.xokv4hb.mongodb.net`
- Qdrant: `194.233.82.241:6333` (no auth)
- Supabase PostgreSQL (user-service)

---

## Prerequisites

```bash
# Tools required
aws --version          # AWS CLI v2
terraform --version    # >= 1.5
kubectl version --client
helm version           # >= 3.x

# Authenticate
aws configure
# Region: ap-southeast-1, Output: json

aws sts get-caller-identity   # verify
```

---

## Step 1 — Deploy VPC + EKS Cluster

```bash
cd iac/cluster

terraform init
terraform apply
# ~15 min — creates VPC, subnets, NAT gateway, EKS 1.33 cluster, node group (t3.medium × 2)
```

**Notes:**
- Kubernetes version is pinned to `1.33` (standard support until Jul 29, 2026)
- Upgrade policy is set to `STANDARD` (auto-upgrades instead of entering paid extended support)
- EKS add-ons (vpc-cni, coredns, kube-proxy) are managed by Terraform

**To upgrade Kubernetes version in the future (one minor version at a time):**
```bash
terraform apply -var="kubernetes_version=1.34"
# wait ~10 min, then:
terraform apply -var="kubernetes_version=1.35"
```

---

## Step 2 — Configure kubectl

```bash
aws eks update-kubeconfig \
  --name e-commerce-eks \
  --region ap-southeast-1

kubectl get nodes   # should show 2 nodes Ready
```

---

## Step 3 — Deploy Traefik (Ingress Controller)

Traefik reads subnet IDs automatically from the cluster Terraform state — no manual subnet IDs needed.

```bash
cd ../traefik

terraform init
terraform apply
# ~3 min — deploys Traefik Helm chart (v41.0.0) and NLB

# Get the NLB hostname (needed for Step 7)
kubectl get svc -n traefik traefik -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'
```

---

## Step 4 — Apply Traefik Middlewares

Must be applied manually after Traefik CRDs are installed:

```bash
cd ../traefik

# Apply CORS middleware (contains allowed origins including CloudFront domain)
kubectl apply -f middlewares-cors-rendered.yaml

# Apply forward-auth middleware (routes auth through user-service)
kubectl apply -f middlewares-auth.yaml

kubectl get middleware -n default   # verify both appear
```

> **To update CORS origins** (e.g. after CloudFront domain changes):
> ```bash
> terraform apply -var="cloudfront_domain=dNEWDOMAIN.cloudfront.net"
> kubectl apply -f middlewares-cors-rendered.yaml
> ```

---

## Step 5 — Deploy S3 + CloudFront (Frontend Hosting)

```bash
cd ../frontend

terraform init
terraform apply
# ~5 min — creates S3 bucket + CloudFront distribution

terraform output cloudfront_domain_name   # note this domain
```

Current CloudFront domain: `d11unvcwv8df4t.cloudfront.net`
CloudFront distribution ID: `E3630N1HELNAEF`

**CloudFront behaviors:**
| Path | Backend | Notes |
|---|---|---|
| `/socket.io/*` | NLB HTTP:80 | WebSocket support, all headers forwarded |
| `/api/*` | NLB HTTP:80 | No caching, auth headers forwarded |
| `/*` (default) | S3 | Static Vue.js SPA |

> **If the NLB hostname changes** (e.g. after cluster destroy+recreate), update `nlb_hostname` in `iac/frontend/terraform.tfvars` and re-apply.

---

## Step 6 — Deploy Application Services (Helm)

All secret values live in `helm/<service>/values-secret.yaml`.

```bash
cd ../../helm

# Infrastructure first
helm upgrade --install rabbitmq ./rabbitmq

# Application services
helm upgrade --install user-service ./user-service -f ./user-service/values-secret.yaml
helm upgrade --install product-service ./product-service -f ./product-service/values-secret.yaml
helm upgrade --install order-service ./order-service -f ./order-service/values-secret.yaml
helm upgrade --install notification-service ./notification-service -f ./notification-service/values-secret.yaml
helm upgrade --install ai-service ./ai-service -f ./ai-service/values-secret.yaml

# Verify all pods are Running
kubectl get pods -n default
```

**Services deployed:**
| Service | Image | Port | Ingress path |
|---|---|---|---|
| user-service | `ghcr.io/phamdatzx/user-service:latest` | 8085 | `/api/user` |
| product-service | `ghcr.io/phamdatzx/product-service:latest` | 8085 | `/api/product` |
| order-service | `ghcr.io/phamdatzx/order-service:latest` | 8085 | `/api/order` |
| notification-service | `ghcr.io/phamdatzx/notification-service:latest` | 8085 | `/api/notification`, `/socket.io` |
| ai-service (API) | `ghcr.io/phamdatzx/ai-service:latest` | 8000 | `/api/ai` |
| ai-service (product-embed-worker) | same image | — | RabbitMQ consumer |
| ai-service (user-vector-worker) | same image | — | RabbitMQ consumer |
| ai-service (cf-worker) | same image | — | CronJob `0 0 * * *` UTC |
| rabbitmq | `rabbitmq:3-management` | 5672/15672 | internal only |

---

## Step 7 — Build & Deploy Frontend

```bash
cd ../frontend

# Install dependencies and build
npm install
npm run build

# Upload built files to S3
aws s3 sync ./dist s3://e-commerce-frontend-web --delete

# Invalidate CloudFront cache
aws cloudfront create-invalidation \
  --distribution-id E3630N1HELNAEF \
  --paths "/*"
```

**Frontend API config (`frontend/.env`):**
```
VITE_BE_API_URL=https://d11unvcwv8df4t.cloudfront.net/api
VITE_USER_API_URL=https://d11unvcwv8df4t.cloudfront.net/api/user/public
...
```
All API calls go through CloudFront (same origin) — no CORS configuration needed in the browser.

---

## Teardown (Destroy All Resources)

```bash
# Remove Helm releases first
helm uninstall ai-service user-service product-service order-service notification-service rabbitmq

# Then destroy infrastructure in reverse order
cd iac/frontend  && terraform destroy
cd ../traefik    && terraform destroy
cd ../cluster    && terraform destroy
```

> Empty the S3 bucket before destroying frontend:
> ```bash
> aws s3 rm s3://e-commerce-frontend-web --recursive
> ```

---

## Re-deploy After Full Teardown

```
1.  cd iac/cluster   → terraform apply
2.  aws eks update-kubeconfig --name e-commerce-eks --region ap-southeast-1
3.  cd iac/traefik   → terraform apply
4.  kubectl apply -f iac/traefik/middlewares-cors-rendered.yaml
5.  kubectl apply -f iac/traefik/middlewares-auth.yaml
6.  cd iac/frontend  → terraform apply
    → note new CloudFront domain
7.  Update cloudfront_domain in iac/traefik terraform.tfvars if domain changed
    → terraform apply -var="cloudfront_domain=<new-domain>"
    → kubectl apply -f iac/traefik/middlewares-cors-rendered.yaml
8.  Update nlb_hostname in iac/frontend/terraform.tfvars if NLB hostname changed
    → cd iac/frontend && terraform apply
9.  Update CLIENT_URL in helm/*/values-secret.yaml with new CloudFront domain
10. cd helm && helm upgrade --install rabbitmq ./rabbitmq
11. helm upgrade --install <each-service> ./<service> -f ./<service>/values-secret.yaml
12. Rebuild frontend with new .env → upload to S3 → invalidate CloudFront
```

---

## Quick Reference — Useful Commands

```bash
# Check all pods
kubectl get pods -n default

# Check ingress routes
kubectl get ingress -n default

# View logs
kubectl logs deployment/user-service
kubectl logs deployment/ai-service
kubectl logs deployment/ai-service-product-embed-worker
kubectl logs deployment/ai-service-user-vector-worker

# Check middlewares
kubectl get middleware -n default

# Get NLB hostname
kubectl get svc -n traefik traefik -o jsonpath='{.status.loadBalancer.ingress[0].hostname}'

# Restart a deployment (e.g. after image update)
kubectl rollout restart deployment/product-service

# Force re-pull latest image
kubectl set image deployment/product-service \
  product-service=ghcr.io/phamdatzx/product-service:latest
```
