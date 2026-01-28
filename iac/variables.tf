//for vpc module
variable "project" {
  type = string
}

variable "cidr_block" {
  type        = string
  description = "CIDR block for the VPC"
}

variable "vpc_name" {
  type        = string
  description = "Name of the VPC"
  default     = "e-commerce-vpc"
}

variable "private_subnet_cidrs" {
  type        = list(string)
  description = "CIDR blocks for private subnets"
}

variable "public_subnet_cidrs" {
  type        = list(string)
  description = "CIDR blocks for public subnets"
}

variable "availability_zones" {
  type        = list(string)
  description = "Availability zones for subnets"
}
