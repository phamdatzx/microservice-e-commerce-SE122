# Read cluster outputs so subnet IDs are always in sync after terraform apply on iac/cluster
data "terraform_remote_state" "cluster" {
  backend = "local"
  config = {
    path = "${path.module}/../cluster/terraform.tfstate"
  }
}

data "aws_eks_cluster" "this" {
  name = var.cluster_name
}

data "aws_eks_cluster_auth" "this" {
  name = var.cluster_name
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.this.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.this.certificate_authority[0].data)
  token                  = data.aws_eks_cluster_auth.this.token
}

provider "helm" {
  kubernetes = {
    host                   = data.aws_eks_cluster.this.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.this.certificate_authority[0].data)
    token                  = data.aws_eks_cluster_auth.this.token
  }
}

# Metrics Server — required for HPA CPU/memory metrics
resource "helm_release" "metrics_server" {
  name             = "metrics-server"
  repository       = "https://kubernetes-sigs.github.io/metrics-server/"
  chart            = "metrics-server"
  namespace        = "kube-system"
  create_namespace = false

  set = [{
    name  = "args[0]"
    value = "--kubelet-insecure-tls"
  }]
}

resource "helm_release" "traefik" {
  name             = "traefik"
  repository       = "https://traefik.github.io/charts"
  chart            = "traefik"
  version          = "41.0.0"
  namespace        = "traefik"
  create_namespace = true

  values = [
    templatefile("${path.module}/values.yaml", {
      subnet_ids = join(",", data.terraform_remote_state.cluster.outputs.public_subnet_ids)
    })
  ]
}

# Render middlewares-cors.yaml with the current CloudFront domain so kubectl apply picks up the right origin.
# After Traefik is deployed:
#   terraform output -raw rendered_cors_middleware > /tmp/middlewares-cors-rendered.yaml
#   kubectl apply -f /tmp/middlewares-cors-rendered.yaml
#   kubectl apply -f iac/traefik/middlewares-auth.yaml
resource "local_file" "middlewares_cors_rendered" {
  filename = "${path.module}/middlewares-cors-rendered.yaml"
  content = templatefile("${path.module}/middlewares-cors.yaml", {
    cloudfront_domain = var.cloudfront_domain
  })
}

# NOTE: These middleware resources are commented out because they depend on Traefik CRDs
# The CRDs are only installed after the Helm chart deploys, causing a chicken-and-egg problem
#
# After Traefik is deployed, apply middlewares manually:
#   kubectl apply -f iac/traefik/middlewares-cors.yaml
#   kubectl apply -f iac/traefik/middlewares-auth.yaml

# resource "kubernetes_manifest" "middleware-cors" {
#   depends_on = [helm_release.traefik]
#
#   manifest = yamldecode(
#     file("${path.module}/middlewares-cors.yaml")
#   )
# }
#
# resource "kubernetes_manifest" "middleware-auth" {
#   depends_on = [helm_release.traefik]
#
#   manifest = yamldecode(
#     file("${path.module}/middlewares-auth.yaml")
#   )
# }
