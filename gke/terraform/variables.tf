variable "gcp_region" {
  description = "Region to launch servers."
  default     = "us-west-1"
}

variable "public_key_path" {
  description = "public key path"
  default = "~/.ssh/aws_key.pub" 
}

variable "db_password" {
  description = "Database password"
}
