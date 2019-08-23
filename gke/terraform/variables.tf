variable "gcp_region" {
  description = "Region to launch servers."
  default     = "us-west-1"
}

variable "public_key_path" {
  description = "public key path"
  default = "~/.ssh/aws_key.pub" 
}

variable "amis_addr" {
  description = "Address server machine image"
  default = {
    us-west-1 = "" // Packer-generated AMI
  }
}

variable "amis_shorten" {
  description = "URL shortener machine images"
  // ami-056ee704806822732" // Unmodifed Amazon AMI (us-west-1)
  default = {
    us-west-1 = "ami-09342686b51a7152f" // Packer-generated AMI
  }
}

variable "amis_expand" {
  description = "URL expander machine images"
  default = {
    us-west-1 = "ami-0d98e0c1037ad9071" // Packer-generated AMI
  }
}



