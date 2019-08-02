// Addr database; URL database shards
// $ export TF_VAR_password=passwd

// Example: www.terraform.io/docs/providers/aws/r/db_instance.html
//   Update parameters within DB parameter group, changes apply to all parameter group DB instances

resource "aws_db_instance" "db_shard0" {
  name                    = "db_shard0"
  allocated_storage       = 20 # GB
  // backup_retention_period  = 7   # days
  // db_subnet_group_name     = "${var.**}"
  // db_subnet_group_name     = "${aws_db_subnet_group.db_subnet_group.id}"
  engine                  = "postgres"
  engine_version          = "9.5.4"
  instance_class          = "db.t2.micro"
  // parameter_group_name    = "url_param_group" # NOTE: Not defined yet
  password = "${var.db_password}"   // Password is stored in TF state file. 
                                    //  Store state file encrypted.
  port                    = 5433
  publicly_accessible     = true
  // storage_encrypted    = {Encryption not supported on db.t2.micro}
  skip_final_snapshot     = true
  storage_type            = "gp2"
  username = "postgres"
  vpc_security_group_ids   = ["${aws_security_group.db.id}"]
}


resource "aws_db_instance" "db_shard1" {
  name                    = "db_shard1"
  allocated_storage       = 20 # GB
  // backup_retention_period  = 7   # days
  // db_subnet_group_name     = "${var.**}"
  // db_subnet_group_name     = "${aws_db_subnet_group.db_subnet_group.id}"
  engine                  = "postgres"
  engine_version          = "9.5.4"
  instance_class          = "db.t2.micro"
  // parameter_group_name    = "url_param_group" # NOTE: Not defined yet
  password = "${var.db_password}"   // Password is stored in TF state file. 
  //  Store state file encrypted.
  port                    = 5433
  publicly_accessible     = true
  skip_final_snapshot     = true
  storage_type            = "gp2"
  username = "postgres"
  vpc_security_group_ids   = ["${aws_security_group.db.id}"]
}


output "url_dbUrl0" {
  value = "${aws_db_instance.db_shard0.address}"
}

output "url_dbUrl1" {
  value = "${aws_db_instance.db_shard1.address}"
}
