# Specify the provider and access details
provider "aws" {
  profile = "default"
  region = "${var.aws_region}"
  version = "~> 2.20"
}

resource "aws_key_pair" "auth" {
  key_name   = "${var.key_name}"
  public_key = "${file(var.public_key_path)}"
}

resource "aws_instance" "example" {
  // ami           = "ami-056ee704806822732" // Unmodifed Amazon AMI
  ami = "ami-0d729a60"
  instance_type = "t2.micro"
  key_name = "${var.key_name}"
  
  tags = {
    Name = "Terraform book"
  }
}

resource "aws_security_group" "instance" {
  name = "terraform-example-instance"

  // TEMP - Web Access
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  // HTTP
  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  // SSH access from anywhere
    ingress {
      from_port   = 22
      to_port     = 22
      protocol    = "tcp"
     cidr_blocks = ["0.0.0.0/0"]
  }

  // outbound internet access
    egress {
      from_port   = 0
      to_port     = 0
      protocol    = "-1"
      cidr_blocks = ["0.0.0.0/0"]
    }
}

}

