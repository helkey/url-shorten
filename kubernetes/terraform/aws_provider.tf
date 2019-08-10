# Specify the provider and access details
provider "aws" {
  profile = "default"
  region = "${var.aws_region}"
  // version = "~> 2.20"
}

resource "aws_vpc" "main" {
  cidr_block = "10.240.0.0/16"
  // enable_dns_support = true (default val)
  enable_dns_hostnames = true
  tags = {
    Name = "kubenetes"
  }
}

resource "aws_vpc_dhcp_options" "dns_resolver" {
  domain_name          = "k8.internal"
  domain_name_servers  = ["AmazonProvidedDNS"]
  tags = {
    Name = "kubenetes"
  }
}

resource "aws_subnet" "k8" {
  vpc_id = "${aws_vpc.main.id}"
  cidr_block = "10.240.0.0/24"
  // availability_zone = ""
  tags = {
    Name = "kubernetes"
  }
}

resource "aws_internet_gateway" "k8gw" {
  vpc_id = "${aws_vpc.main.id}"
  tags = {
    Name = "kubernetes"
  }
}

resource "aws_route_table" "r" {
  vpc_id = "${aws_vpc.main.id}"
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.k8gw.id}"
  }
  tags = {
    Name = "kubernetes"
  }
}

resource "aws_security_group" "k8" {
  name = "Kubernetes security group"

  // TEMP - Web Access
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress { // Any trafficx
    from_port   = 0
    to_port     = 0
    protocol    = "-1" // all
    cidr_blocks = ["10.240.0.0/16"]
  }
  ingress { // SSH access from anywhere
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  
  ingress { // Kubernetes SSL
    from_port   = 6443
    to_port     = 6443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress { // Any-any within security group
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    self = true
    // source_security_group_id="${aws_security_group.kubernetes.id}"
  }
  
  
  tags = {
    Name = "kubernetes"
  }
}

resource "aws_key_pair" "auth" {
  key_name   = "${var.key_name}"
  public_key = "${file(var.public_key_path)}"
}



// When using default VPC, AWS provides default subnets in two zones
//    as required for high-availability RDS
// resource "aws_default_vpc" "default" {
// }


