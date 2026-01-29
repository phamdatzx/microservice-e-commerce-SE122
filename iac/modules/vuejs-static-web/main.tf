# S3 Bucket for static website hosting
resource "aws_s3_bucket" "vuejs_web" {
  bucket = var.bucket_name

  tags = {
    Name    = "${var.project}-vuejs-web"
    Project = var.project
  }
}

# S3 Bucket Website Configuration
resource "aws_s3_bucket_website_configuration" "vuejs_web" {
  bucket = aws_s3_bucket.vuejs_web.id

  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "index.html"  # For SPA routing
  }
}

# S3 Bucket Public Access Block (we'll use CloudFront OAI instead)
resource "aws_s3_bucket_public_access_block" "vuejs_web" {
  bucket = aws_s3_bucket.vuejs_web.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# CloudFront Origin Access Identity
resource "aws_cloudfront_origin_access_identity" "vuejs_web" {
  comment = "OAI for ${var.project} Vue.js web"
}

# S3 Bucket Policy for CloudFront
resource "aws_s3_bucket_policy" "vuejs_web" {
  bucket = aws_s3_bucket.vuejs_web.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowCloudFrontAccess"
        Effect = "Allow"
        Principal = {
          AWS = aws_cloudfront_origin_access_identity.vuejs_web.iam_arn
        }
        Action   = "s3:GetObject"
        Resource = "${aws_s3_bucket.vuejs_web.arn}/*"
      }
    ]
  })
}

# CloudFront Distribution
resource "aws_cloudfront_distribution" "vuejs_web" {
  enabled             = true
  is_ipv6_enabled     = true
  default_root_object = "index.html"
  comment             = "${var.project} Vue.js static web"
  price_class         = var.cloudfront_price_class

  origin {
    domain_name = aws_s3_bucket.vuejs_web.bucket_regional_domain_name
    origin_id   = "S3-${aws_s3_bucket.vuejs_web.id}"

    s3_origin_config {
      origin_access_identity = aws_cloudfront_origin_access_identity.vuejs_web.cloudfront_access_identity_path
    }
  }

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "S3-${aws_s3_bucket.vuejs_web.id}"

    forwarded_values {
      query_string = false
      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "redirect-to-https"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
    compress               = true
  }

  # Custom error response for SPA routing
  custom_error_response {
    error_code         = 404
    response_code      = 200
    response_page_path = "/index.html"
  }

  custom_error_response {
    error_code         = 403
    response_code      = 200
    response_page_path = "/index.html"
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = var.custom_domain == "" ? true : false
    acm_certificate_arn            = var.acm_certificate_arn
    ssl_support_method             = var.custom_domain == "" ? null : "sni-only"
    minimum_protocol_version       = var.custom_domain == "" ? null : "TLSv1.2_2021"
  }

  dynamic "aliases" {
    for_each = var.custom_domain != "" ? [var.custom_domain] : []
    content {
      aliases = [var.custom_domain]
    }
  }

  tags = {
    Name    = "${var.project}-vuejs-web-cdn"
    Project = var.project
  }
}

# Optional: Route 53 record (if custom domain is provided)
resource "aws_route53_record" "vuejs_web" {
  count   = var.custom_domain != "" && var.route53_zone_id != "" ? 1 : 0
  zone_id = var.route53_zone_id
  name    = var.custom_domain
  type    = "A"

  alias {
    name                   = aws_cloudfront_distribution.vuejs_web.domain_name
    zone_id                = aws_cloudfront_distribution.vuejs_web.hosted_zone_id
    evaluate_target_health = false
  }
}
