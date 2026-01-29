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
  depends_on = [ module.vpc ]
  source = "./modules/eks"  

  project = var.project
  subnet_ids = module.vpc.private_subnet_ids
}