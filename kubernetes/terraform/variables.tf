variable "aws_region" {
  description = "AWS region to launch servers."
  default     = "us-west-1"
}

variable "db_password" {
  description = "Password for addr and url databases"
}

variable "key_name" {
  default = "aws_key"
}

variable "public_key_path" {
  description = "public key path"
  default = "~/.ssh/aws_key.pub" 
}

variable "amis_k8" {
  description = "Amazon instance machine images"
  default = {
    us-west-1 = ""
  }
}

