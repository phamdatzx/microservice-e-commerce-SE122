output "vpc_id" {
  value       = aws_vpc.main.id
  description = "ID of the VPC"
}

output "vpc_cidr_block" {
  value       = aws_vpc.main.cidr_block
  description = "CIDR block of the VPC"
}

output "private_subnet_ids" {
  value       = aws_subnet.private[*].id
  description = "IDs of private subnets"
}

output "public_subnet_ids" {
  value       = aws_subnet.public[*].id
  description = "IDs of public subnets"
}

output "public_subnet_cidr_blocks" {
  value       = aws_subnet.public[*].cidr_block
  description = "CIDR blocks of public subnets"
}