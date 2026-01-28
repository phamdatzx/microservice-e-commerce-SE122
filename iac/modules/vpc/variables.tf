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

variable "private_subnet_cidr" {
  type        = string
  description = "CIDR block for private subnet"
}

variable "public_subnet_cidr" {
  type        = string
  description = "CIDR block for public subnet"
}

variable "availability_zone" {
  type        = string
  description = "Availability zones for subnets"
}
