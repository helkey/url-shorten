
# Security group so ELB accessible from web
resource "aws_security_group" "elb" {
  name        = "security_elb"
  vpc_id      = "${aws_vpc.default.id}"

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
  vpc_id      = "${aws_vpc.default.id}"

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

resource "aws_elb" "web" {
  name = "aws-elb"

  subnets         = ["${aws_subnet.default.id}"]
  security_groups = ["${aws_security_group.elb.id}"]
  instances       = ["${aws_instance.web.id}"]

  listener {
    instance_port     = 80
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }
}

resource "aws_instance" "web" {
  # Connection block tells provisioner how communicate with instances
  connection {
    # Connection will use local SSH agent for authentication.
    # The default username for AMI
    user = "ubuntu"
    host = "${self.public_ip}"
  }

  instance_type = "t2.micro"

  # Lookup the AMI based on region
  ami = "${lookup(var.aws_amis, var.aws_region)}"

  # Name of SSH keypair created above.
  key_name = "${aws_key_pair.auth.id}"

  # Security group to allow HTTP and SSH access
  vpc_security_group_ids = ["${aws_security_group.default.id}"]

  # Launch into same subnet as ELB for testing.
  # Later switch to separate private subnet for backend instances.
  subnet_id = "${aws_subnet.default.id}"
}
