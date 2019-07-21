// URL database shards
// $ export TF_VAR_password=passwd

// Example: www.terraform.io/docs/providers/aws/r/db_instance.html
// Update parameters within DB parameter group, changes apply to all parameter group DB instances
resource "aws_db_instance" "shard0" {
  name                    = "url_shard0"
  allocated_storage       = 20 # GB
  // backup_retention_period  = 7   # days
  // db_subnet_group_name     = "${var.**}"
  // db_subnet_group_name     = "${aws_db_subnet_group.db_subnet_group.id}"
  engine                  = "postgres"
  engine_version          = "9.5.4"
  instance_class          = "db.t2.micro"
  // parameter_group_name    = "url_param_group" # NOTE: Not defined yet

  password = "${var.db_password}"   // Password is stored in TF state file. Store state file encrpted,
                                    //  or modify afterward using AWS 
  port                    = 5433
  publicly_accessible     = true
  // storage_encrypted    = {true} # Encryption not supported on db.t2.micro 
  skip_final_snapshot     = true
  storage_type            = "gp2"
  username = "postgres"
  vpc_security_group_ids   = ["${aws_security_group.db.id}"]
}

resource "aws_security_group" "db" {
  // vpc_id      = "${aws_vpc.default.id}"
  vpc_id      = "${aws_default_vpc.default.id}"
  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks     = ["0.0.0.0/0"] // TODO: restrict CIDR
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
//  subnet_ids = ["${aws_subnet.default.id}"] # NEED to specify subnets from two regions
//
//  tags = {
//    Name = "My DB subnet group"
//  }
//}