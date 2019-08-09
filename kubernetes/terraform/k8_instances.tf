

resource "aws_instance" "controller_0" {
  ami = "${lookup(var.amis_k8, var.aws_region)}"
  // ami = "ami-056ee704806822732" // Unmodifed Amazon AMI (us-west-1)
  instance_type = "t2.micro"
  key_name = "${var.key_name}"
  vpc_security_group_ids = [aws_security_group.k8.id]

  associate_public_ip_address = true
  private_ip = "10.240.0.10"
  subnet_id   = "${aws_subnet.k8.id}"
  tags = {
    Name = "k8_control"
  }
}

