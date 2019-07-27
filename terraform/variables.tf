variable "public_key_path" {
  description = "public key path"
  default = "~/.ssh/terraform_rsa.pub"
}

variable "key_name" {
  description = "key name"
  default = "~/.ssh/authorized_keys/url-key-uswest" // DONT add .pem!!!
}

variable "db_password" {
}

variable "aws_region" {
  description = "AWS region to launch servers."
  default     = "us-west-1"
}

variable "addr_amis" {
  description = "Amazon instance machine images"
  default = {
    us-west-1 = "ami-056ee704806822732" // Unmodifed Amazon AMI
    // not responding: ami-0e29c20cd0c57c5ed, ami-087fb2fd7e9c79e6e
    us-west-2 = "ami-052d9eee6d5a9bf35"
  }
}
