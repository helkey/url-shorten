
resource "aws_instance" "addr_server" {
  # Connection block tells provisioner how communicate with instances
  connection {
    # Connection will use local SSH agent for authentication.
    # The default username for AMI
    user = "ec2-user" //  "ubuntu"
    host = "${self.public_ip}"
  }

  instance_type = "t2.micro"

  # Lookup the AMI based on region
  ami = "${lookup(var.addr_amis, var.aws_region)}"

  # Name of SSH keypair created above.
  key_name = "${aws_key_pair.auth.id}"

  # Security group to allow HTTP and SSH access
  vpc_security_group_ids = ["${aws_security_group.addr.id}"]

  # Launch into same subnet as ELB for testing.
  # Later switch to separate private subnet for backend instances.
  # subnet_id = "${aws_subnet.default.id}"
}


// 
resource "aws_security_group" "addr" {
  // vpc_id      = "${aws_vpc.default.id}"
  vpc_id      = "${aws_default_vpc.default.id}"
  ingress {
    // from_port   = 22 // SSH access
    // to_port     = 22
    // protocol    = "tcp"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks     = ["0.0.0.0/0"] // TODO: restrict CIDR?
}
  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["0.0.0.0/0"]
  }
}

// resource "aws_db_subnet_group" "db_subnet_group" {
//  name       = "main"
//# DBs need subnets from two regions to meet avail requirements
