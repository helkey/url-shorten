
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
  vpc_security_group_ids = ["${aws_security_group.db.id}"]

  # Launch into same subnet as ELB for testing.
  # Later switch to separate private subnet for backend instances.
  # subnet_id = "${aws_subnet.default.id}"
}
