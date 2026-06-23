variable "cluster_name" {
  type        = string
  description = "Name of the EKS cluster"
}

variable "cloudfront_domain" {
  type        = string
  description = "CloudFront distribution domain to allow in CORS (e.g. dxxxxxx.cloudfront.net). Update after iac/frontend is applied."
  default     = ""
}
