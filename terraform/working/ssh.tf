// https://learning.oreilly.com/library/view/terraform-up/9781492046899/ch02.html#deploy_single_server

variable "key_name" {
  default = "aws_public_key"
}

variable "public_key_path" {
  description = "public key path"
  default = "~/.ssh/aws_key.pub"
}


provider "aws" {
  region = "us-west-1"
  version = "~> 2.20"
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

resource "aws_key_pair" "auth" {
  key_name   = "${var.key_name}"
  public_key = "${file(var.public_key_path)}"
}

resource "aws_instance" "example" {
  ami                    = "ami-056ee704806822732" // Unmodifed AMI (us-west-1)
  instance_type          = "t2.micro"
  key_name = "${var.key_name}"
  vpc_security_group_ids = [aws_security_group.instance.id]

  //  Instance 'smoke test'
  user_data = <<-EOF
              #!/bin/bash
              echo "Hello, World" > index.html
              nohup busybox httpd -f -p 8080 &
              EOF

  tags = {
      Name = "ssh-example"
  }
}



