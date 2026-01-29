output "bucket_name" {
  description = "S3 bucket name"
  value       = aws_s3_bucket.vuejs_web.id
}

output "bucket_arn" {
  description = "S3 bucket ARN"
  value       = aws_s3_bucket.vuejs_web.arn
}

output "cloudfront_distribution_id" {
  description = "CloudFront distribution ID"
  value       = aws_cloudfront_distribution.vuejs_web.id
}

output "cloudfront_domain_name" {
  description = "CloudFront distribution domain name"
  value       = aws_cloudfront_distribution.vuejs_web.domain_name
}

output "website_url" {
  description = "Website URL (CloudFront or custom domain)"
  value       = var.custom_domain != "" ? "https://${var.custom_domain}" : "https://${aws_cloudfront_distribution.vuejs_web.domain_name}"
}
