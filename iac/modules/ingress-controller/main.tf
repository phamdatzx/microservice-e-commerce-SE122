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

resource "helm_release" "traefik" {
  name       = "traefik"
  repository = "https://traefik.github.io/charts"
  chart      = "traefik"
  namespace  = "traefik"
  create_namespace = true

  values = [
    templatefile("${path.module}/values.yaml", {
      subnet_ids = join(",", var.subnet_ids)
    })
  ]
}

# NOTE: These middleware resources are commented out because they depend on Traefik CRDs
# The CRDs are only installed after the Helm chart deploys, causing a chicken-and-egg problem
# 
# To apply these middlewares, run kubectl after Traefik is deployed:
#   kubectl apply -f modules/ingress-controller/middlewares-cors.yaml
#   kubectl apply -f modules/ingress-controller/middlewares-auth.yaml

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
