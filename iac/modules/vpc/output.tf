output "vpc_id" {
  value       = aws_vpc.main.id
  description = "ID of the VPC"
}

output "vpc_cidr_block" {
  value       = aws_vpc.main.cidr_block
  description = "CIDR block of the VPC"
}

output "private_subnet_id" {
  value       = aws_subnet.private.id
  description = "ID of private subnet"
}

output "public_subnet_id" {
  value       = aws_subnet.public.id
  description = "ID of public subnet"
}

output "public_subnet_cidr_block" {
  value       = aws_subnet.public.cidr_block
  description = "CIDR block of public subnet"
}