
# Security group so ELB accessible from web
resource "aws_security_group" "elb" {
  name        = "security_elb"
  // vpc_id      = "${aws_vpc.default.id}"

  # HTTP access from anywhere
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Security group to access instances over SSH / HTTP
resource "aws_security_group" "default" {
  name        = "security_default"
  // vpc_id      = "${aws_vpc.default.id}"

  # SSH access from anywhere
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # HTTP access from the VPC
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }

  # outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "shorten_server" {
  ami = "${lookup(var.amis_shorten, var.aws_region)}"
  instance_type = "t2.micro"
  key_name = "${var.key_name}"
  vpc_security_group_ids = [aws_security_group.instance.id]

  tags = {
    Name = "Shorten server"
  }
}

resource "aws_instance" "expand_server" {
  ami = "${lookup(var.amis_expand, var.aws_region)}"
  instance_type = "t2.micro"
  key_name = "${var.key_name}"
  vpc_security_group_ids = [aws_security_group.instance.id]

  tags = {
    Name = "Expand server"
  }
}

  
output "shorten_ip" {
  value="${aws_instance.shorten_server.public_ip}"
}

output "expand_ip" {
  value="${aws_instance.expand_server.public_ip}"
}
