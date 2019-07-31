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

variable "addr_amis" {
  description = "Amazon instance machine images"
  default = {
    us-west-1 = "ami-065033204f8aad05a" // Packer AMI
  }
}

