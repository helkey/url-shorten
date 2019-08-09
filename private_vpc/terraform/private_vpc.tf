variable "private_cidr" {
  description = "CIDR for the Private Subnet"
  default = "10.0.1.0/24"
}


resource "aws_subnet" "private" {
  vpc_id = "${aws_vpc.default.id}"
  cidr_block = "${var.private_cidr}"
  availability_zone = ""

  tags {
    Name = "private"
  }
}

resource "aws_route_table" "private" {
  vpc_id = "${aws_vpc.default.id}"
  route {
    cidr_block = "0.0.0.0/0"
    instance_id = "${aws_instance.nat.id}"
  }
}

resource "aws_route_table_association" "private" {
  subnet_id = "${aws_subnet.private.id}"
  route_table_id = "${aws_route_table.private.id}"
}

resource "aws_security_group" "db" {
  vpc_id = "${aws_vpc.default.id}"
  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["${var.private_cidr}"]
  }
  ingress { # DB
    from_port = 5433
    to_port = 5433
    protocol = "tcp"
    security_groups = ["${aws_security_group.public.id}"]
  }
  ingress { # DB error messages
    from_port = -1
    to_port = -1
    protocol = "icmp"
    cidr_blocks = ["${var.private_cidr}"]
  }

  egress { # HTTP to internet
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  egress { # HTTPS to internet
    from_port = 443
    to_port = 443
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "database"
  }
}

resource "aws_instance" "db-1" {

}
