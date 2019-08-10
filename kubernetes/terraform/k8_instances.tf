// 

# IF you create the role in Terraform
resource "aws_iam_role" "k8" {
  name = "k8_role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
  tags = {
      tag-key = "iam_role_k8"
  }
}

resource "aws_iam_role_policy" "k8" {
  name = "k8_policy"
  role = "${aws_iam_role.k8.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}


// Provide an IAM instance profile: https://www.terraform.io/docs/providers/aws/r/iam_instance_profile.html
data "aws_iam_instance_profile" "k8" {
  name = "name_of_instance_profile" 
}


//-- KUBERNETES CONTROLLER  --
resource "aws_instance" "controller_0" {
  ami = "${lookup(var.amis_k8, var.aws_region)}"
  instance_type = "t2.micro"
  associate_public_ip_address = true
  iam_instance_profile = "kubernetes"
  key_name = "${var.key_name}"
  private_ip = "10.240.0.10"
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
  iam_instance_profile = "kubernetes"
  key_name = "${var.key_name}"
  private_ip = "10.240.0.20"
  subnet_id   = "${aws_subnet.k8.id}"
  vpc_security_group_ids = [aws_security_group.k8.id]
  tags = {
    Name = "k8_control"
  }
}



