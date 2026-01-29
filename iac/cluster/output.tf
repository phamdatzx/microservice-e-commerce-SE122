output "cluster_name" {
  value = module.eks.cluster_name
}

output "private_subnet_ids" {
  value = module.vpc.private_subnet_ids
}