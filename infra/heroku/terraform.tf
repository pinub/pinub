terraform {
  backend "s3" {
    encrypt = true
    bucket  = "tf-infra"
    key     = "pinub-heroku.tfstate"
    region  = "eu-central-1"
  }
}

variable "domain" {
  type    = "string"
  default = "pinub.com"
}

variable "heroku_email" {}
variable "heroku_api_key" {}
variable "heroku_webhook_url_sentry" {}
variable "cloudflare_email" {}
variable "cloudflare_token" {}

# HEROKU

provider "heroku" {
  email   = "${var.heroku_email}"
  api_key = "${var.heroku_api_key}"
}

resource "heroku_app" "pinub" {
  name   = "pinub"
  region = "eu"
  stack  = "heroku-16"

  config_vars = {
    CSS = "https://${aws_cloudfront_distribution.assets.domain_name}/styles.css?${substr(aws_s3_bucket_object.css.etag, 0, 5)}"
    JS  = "https://${aws_cloudfront_distribution.assets.domain_name}/basic.js?${substr(aws_s3_bucket_object.js.etag, 0, 5)}"
  }
}

resource "heroku_domain" "pinub_com" {
  app      = "${heroku_app.pinub.name}"
  hostname = "${var.domain}"
}

resource "heroku_domain" "www_pinub_com" {
  app      = "${heroku_app.pinub.name}"
  hostname = "www.${var.domain}"
}

resource "heroku_addon" "pinub-database" {
  app  = "${heroku_app.pinub.name}"
  plan = "heroku-postgresql:hobby-dev"
}

resource "heroku_addon" "pinub-webhook-sentry" {
  app  = "${heroku_app.pinub.name}"
  plan = "deployhooks:http"

  config {
    url = "${var.heroku_webhook_url_sentry}"
  }
}

# CLOUDFLARE

provider "cloudflare" {
  email = "${var.cloudflare_email}"
  token = "${var.cloudflare_token}"
}

resource "cloudflare_record" "pinub_com" {
  domain  = "${var.domain}"
  name    = "pinub.com"
  value   = "${heroku_domain.pinub_com.cname}"
  type    = "CNAME"
  proxied = true
}

resource "cloudflare_record" "www_pinub_com" {
  domain  = "${var.domain}"
  name    = "www"
  value   = "${heroku_domain.www_pinub_com.cname}"
  type    = "CNAME"
  proxied = true
}

# AWS

provider "aws" {
  region = "eu-central-1"
}

resource "aws_s3_bucket" "assets" {
  bucket = "pinub-static"
  acl    = "public-read"

  tags {
    name = "pinub"
  }
}

resource "aws_s3_bucket_object" "css" {
  bucket = "${aws_s3_bucket.assets.bucket}"
  key    = "styles.css"
  source = "../../static/css/styles.css"
  etag   = "${md5(file("../../static/css/styles.css"))}"
  acl    = "public-read"

  content_type = "text/css"
}

resource "aws_s3_bucket_object" "js" {
  bucket = "${aws_s3_bucket.assets.bucket}"
  key    = "basic.js"
  source = "../../static/js/basic.js"
  etag   = "${md5(file("../../static/js/basic.js"))}"
  acl    = "public-read"

  content_type = "application/javascript"
}

resource "aws_cloudfront_distribution" "assets" {
  enabled         = true
  is_ipv6_enabled = true
  price_class     = "PriceClass_100"

  origin {
    domain_name = "${aws_s3_bucket.assets.bucket_domain_name}"
    origin_id   = "s3_assets_origin"
  }

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    compress         = true
    target_origin_id = "s3_assets_origin"

    forwarded_values {
      query_string = true

      cookies {
        forward = "none"
      }
    }

    viewer_protocol_policy = "https-only"
    min_ttl                = 0
    default_ttl            = 86400
    max_ttl                = 31536000
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  tags {
    name = "pinub"
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}
