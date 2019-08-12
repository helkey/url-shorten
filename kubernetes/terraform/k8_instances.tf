// k8_instances

// medium.com/@devopslearning/aws-iam-ec2-instance-role-using-terraform-fa2b21488536
# IAM role
resource "aws_iam_role" "k8" {
  name = "k8_role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17", 
  "Statement": [
    {"Effect": "Allow", "Principal": { "Service": "ec2.amazonaws.com"}, "Action": "sts:AssumeRole"}
  ]
}
EOF
  tags = {
       Name = "iam_role_k8"
  }
}

// IAM instance profile:
//  www.terraform.io/docs/providers/aws/r/iam_instance_profile.html
resource "aws_iam_instance_profile" "k8" {
  name = "k8_profile"
  role = "${aws_iam_role.k8.name}"
}

// Role policy gives EC2 instance full access to S3 storage
resource "aws_iam_role_policy" "k8" {
  name = "k8_policy"
  role = "${aws_iam_role.k8.id}"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {"Action": ["s3:*"], "Effect": "Allow", "Resource": "*"}
  ]
}
EOF
}

resource "aws_iam_policy" "k8" {
  name        = "k8_policy"
  path        = "/"
  description = "Kubernetes policy"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {"Effect": "Allow", "Action": ["ec2:*"], "Resource": ["*"]},
    {"Effect": "Allow", "Action": ["elasticloadbalancing:*"], "Resource": ["*"]},
    {"Effect": "Allow", "Action": ["route53:*"], "Resource": ["*"]},
    {"Effect": "Allow", "Action": ["ecr:*"], "Resource": "*"}
  ]
}
EOF
}


//-- KUBERNETES CONTROLLER  --
resource "aws_instance" "controller_0" {
  ami = "${lookup(var.amis_k8, var.aws_region)}"
  instance_type = "t2.micro"
  associate_public_ip_address = true
  iam_instance_profile = "k8_profile"
  key_name = "${var.key_name}"
  private_ip = "10.240.0.10"
  source_dest_check = false // enable traffic between foreign subnets
  subnet_id   = "${aws_subnet.k8.id}"
  vpc_security_group_ids = [aws_security_group.k8.id]
  tags = {
    Name = "k8_control"
  }
}

//-- KUBERNETES WORKER  --
resource "aws_instance" "worker_0" {
  ami = "${lookup(var.amis_k8, var.aws_region)}"
  instance_type = "t2.micro"
  associate_public_ip_address = true
  iam_instance_profile = "k8_profile"
  key_name = "${var.key_name}"
  private_ip = "10.240.0.20"
  source_dest_check = false
  subnet_id   = "${aws_subnet.k8.id}"
  vpc_security_group_ids = [aws_security_group.k8.id]
  tags = {
    Name = "k8_control"
  }
}



