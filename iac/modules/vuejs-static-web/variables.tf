variable "project" {
  description = "Project name"
  type        = string
}

variable "bucket_name" {
  description = "S3 bucket name for static website hosting"
  type        = string
}

variable "cloudfront_price_class" {
  description = "CloudFront distribution price class"
  type        = string
  default     = "PriceClass_100"  # Use only North America and Europe
}

variable "custom_domain" {
  description = "Custom domain name (e.g., app.example.com). Leave empty to use CloudFront domain"
  type        = string
  default     = ""
}

variable "acm_certificate_arn" {
  description = "ACM certificate ARN for custom domain (must be in us-east-1 region for CloudFront)"
  type        = string
  default     = ""
}

variable "route53_zone_id" {
  description = "Route 53 hosted zone ID for custom domain"
  type        = string
  default     = ""
}
