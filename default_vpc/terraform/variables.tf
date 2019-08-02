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

variable "amis_addr" {
  description = "Amazon instance machine images"
  default = {
    us-west-1 = "ami-044d49c65c0c6abf2" // Packer-generated AMI
  }
}

variable "amis_shorten" {
  description = "URL shortener machine images"
  // ami-056ee704806822732" // Unmodifed Amazon AMI (us-west-1)
  default = {
    us-west-1 = "ami-09342686b51a7152f" // Unmodifed Amazon AMI (us-west-1)
  }
}

variable "amis_expand" {
  description = "URL expander machine images"
  default = {
    us-west-1 = "" // Packer-generated AMI
  }
}



