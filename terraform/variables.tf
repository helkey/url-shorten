variable "public_key_path" {
  description = "public key path"
  default = "~/.ssh/terraform_rsa.pub"
}

variable "key_name" {
  description = "key name"
  default = "~/.ssh/authorized_keys/url-key-uswest.pem"
}

variable "aws_region" {
  description = "AWS region to launch servers."
  default     = "us-west-2"
}

variable "aws_amis" {
  description = "Amazon instance machine images"
  default = {
    us-east-1 = "ami-1d4e7a66"
    us-west-1 = "ami-969ab1f6"
    us-west-2 = "ami-8803e0f0"
  }
}
