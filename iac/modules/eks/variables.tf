variable "project" {
  type        = string
  description = "Name of the project"
}

variable "subnet_ids" {
  type        = list(string)
  description = "List of subnet IDs for the EKS cluster"
}