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

resource "aws_security_group" "instance" {
  name = "terraform-example-instance"

  // TEMP - Web Access
  ingress {
    from_port   = 80
    to_port     = 80
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

  // Database access
  ingress {
    from_port   = 5433
    to_port     = 5433
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  // Internal HTTP services
  ingress {
    from_port   = 8088
    to_port     = 8088
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


// When using default VPC, AWS provides default subnets in two zones
//    as required for high-availability RDS
resource "aws_default_vpc" "default" {
  tags = {
    Name = "Default VPC"
  }
}


resource "aws_instance" "addr_server" {
  ami = "${lookup(var.addr_amis, var.aws_region)}"
  // ami = "ami-056ee704806822732" // Unmodifed Amazon AMI (us-west-1)
  instance_type = "t2.micro"
  key_name = "${var.key_name}"
  vpc_security_group_ids = [aws_security_group.instance.id]

  //  Instance 'smoke test'
  //user_data = <<-EOF
  //            #!/bin/bash
  //            echo "Hello, World" > index.html
  //            nohup busybox httpd -f -p 8080 &
  //            EOF
  
  tags = {
    Name = "Unmodifed AMI ...2732"
  }
}

