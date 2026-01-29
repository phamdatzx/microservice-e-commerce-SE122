module "vpc" {
  source = "./modules/vpc"

  project = var.project
  cidr_block = var.cidr_block
  vpc_name = var.vpc_name
  private_subnet_cidrs = var.private_subnet_cidrs
  public_subnet_cidrs = var.public_subnet_cidrs
  availability_zones = var.availability_zones
}

module "eks" {
  source = "./modules/eks"  

  project = var.project
  subnet_ids = module.vpc.private_subnet_ids
}

#Temporarily commented out - will enable after EKS is created
module "ingress-controller" {
  source = "./modules/ingress-controller"

  cluster_name    = module.eks.cluster_name
  subnet_ids = module.vpc.public_subnet_ids
}

# Vue.js Static Web - CloudFront + S3
module "vuejs-static-web" {
  source = "./modules/vuejs-static-web"

  project     = var.project
  bucket_name = "${var.project}-frontend-web"
}
