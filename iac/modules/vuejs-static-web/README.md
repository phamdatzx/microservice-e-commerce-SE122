# Vue.js Static Web Module

This Terraform module deploys a Vue.js (or any static) web application to AWS using S3 and CloudFront.

## Features

- **S3 Bucket**: Hosts static files with website configuration
- **CloudFront CDN**: Global content delivery with caching
- **HTTPS**: Enforced with redirect from HTTP
- **SPA Support**: Proper error handling for Vue Router (client-side routing)
- **Custom Domain** (Optional): Support for custom domains with ACM certificates
- **Route 53** (Optional): Automatic DNS record creation

## Usage

### Basic Usage (CloudFront Domain)

```hcl
module "vuejs_web" {
  source = "./modules/vuejs-static-web"

  project     = var.project
  bucket_name = "my-app-frontend-bucket"
}
```

### With Custom Domain

```hcl
module "vuejs_web" {
  source = "./modules/vuejs-static-web"

  project              = var.project
  bucket_name          = "my-app-frontend-bucket"
  custom_domain        = "app.example.com"
  acm_certificate_arn  = "arn:aws:acm:us-east-1:123456789:certificate/abc-123"
  route53_zone_id      = "Z1234567890ABC"
}
```

## Prerequisites

- **ACM Certificate**: If using a custom domain, create an ACM certificate in **us-east-1** region (required for CloudFront)
- **Route 53 Hosted Zone**: If you want automatic DNS record creation

## Deploying Your Vue.js App

After Terraform creates the infrastructure:

### 1. Build your Vue.js app

```bash
cd /path/to/your/vuejs/app
npm run build
```

### 2. Upload to S3

```bash
aws s3 sync ./dist s3://YOUR-BUCKET-NAME --delete
```

### 3. Invalidate CloudFront cache (if needed)

```bash
aws cloudfront create-invalidation \
  --distribution-id YOUR-DISTRIBUTION-ID \
  --paths "/*"
```

## Outputs

- `bucket_name`: S3 bucket name
- `bucket_arn`: S3 bucket ARN
- `cloudfront_distribution_id`: CloudFront distribution ID
- `cloudfront_domain_name`: CloudFront domain name
- `website_url`: Full website URL

## Notes

- The module configures the S3 bucket to redirect 404 and 403 errors to `/index.html` for proper SPA routing
- CloudFront is configured to compress content automatically
- Public access to the S3 bucket is blocked; access is only through CloudFront
- Default cache TTL is 1 hour (3600 seconds)
