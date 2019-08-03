variable "public_cidr" {
  default = "10.0.0.0/24"
}

resource "aws_route_table" "public" {
  vpc_id = "${aws_vpc.default.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.default.id}"
  }

  tags {
    Name = "Public Subnet"
  }
}

resource "aws_route_table_association" "public" {
  subnet_id = "${aws_subnet.public.id}"
  route_table_id = "${aws_route_table.public.id}"
}

resource "aws_subnet" "public" {
  vpc_id = "${aws_vpc.default.id}"

  cidr_block = "${var.public_cidr}"
  availability_zone = ""

  tags {
    Name = "Public Subnet"
  }
}

resource "aws_route_table" "eu-west-1a-public" {
  vpc_id = "${aws_vpc.default.id}"

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.default.id}"
  }

  tags {
    Name = "Public Subnet"
  }
}

resource "aws_route_table_association" "public" {
  subnet_id = "${aws_subnet.public.id}"
  route_table_id = "${aws_route_table.public.id}"
}
