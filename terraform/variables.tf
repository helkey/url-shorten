variable "aws_region" {
  description = "AWS region to launch servers."
  default     = "us-west-1"
}

// variable "db_password" {
//}

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
    us-west-1 = "ami-056ee704806822732" // Unmodifed Amazon AMI
    us-west-2 = "ami-052d9eee6d5a9bf35"
  }
}

