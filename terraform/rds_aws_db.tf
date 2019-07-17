// URL database shards

// Example: www.terraform.io/docs/providers/aws/r/db_instance.html
// Update parameters within DB parameter group, changes apply to all parameter group DB instances
resource "aws_db_instance" "shard0" {
  // provider             = "postgresql.shard0"
  name                 = "url_shard0"
  parameter_group_name = "url_param_group"
  instance_class       = "db.t2.micro"
  engine               = "postgres"
  engine_version       = "9.5.4"
  storage_type         = "gp2"
  allocated_storage    = 20
  username = "postgres_url"
  password = "rh1cr2g9"
}






