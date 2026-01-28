# Terraform Apply Instructions

## The Problem

The ingress-controller module contains Kubernetes and Helm providers that need to read the EKS cluster data during planning. However, the cluster doesn't exist yet, causing this error:

```
Error: reading EKS Cluster (e-commerce-eks): couldn't find resource
```

## Solution: Two-Stage Apply

### Stage 1: Create VPC and EKS Cluster

First, create only the VPC and EKS infrastructure:

```bash
terraform apply -target=module.vpc -target=module.eks
```

This will create:
- VPC with public and private subnets
- EKS cluster
- EKS node group
- All IAM roles and policies

**Wait time**: ~15-20 minutes for EKS cluster to become ready

### Stage 2: Deploy Ingress Controller

After the EKS cluster is created and ready, apply the ingress-controller:

```bash
terraform apply
```

This will now succeed because the EKS cluster exists and the data sources can read it.

## Alternative: One-Command Apply (if you comment out ingress-controller first)

If you want to avoid the two-stage process in the future, you could:

1. Comment out the `ingress-controller` module in `main.tf`
2. Run `terraform apply` to create VPC and EKS
3. Uncomment the `ingress-controller` module
4. Run `terraform apply` again

## Why This Happens

The ingress-controller module has **provider configurations** inside it that depend on data sources:

```hcl
data "aws_eks_cluster" "this" {
  name = var.cluster_name  # This cluster doesn't exist yet!
}

provider "kubernetes" {
  host = data.aws_eks_cluster.this.endpoint  # Evaluated during planning
  ...
}
```

Terraform evaluates data sources and providers during the planning phase, but the cluster doesn't exist yet, causing the error.

## Next Steps

Run the Stage 1 command and wait for completion, then run Stage 2.
