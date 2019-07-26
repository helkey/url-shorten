// $terraform output shard_0 shard_1... > "aws.ip.txt"

//output "elb_address" {
// value = "${aws_elb.web.dns_name}"
//}

// output "addr_shard0" {
output "addr_dbAddr" {
  // value = "${aws_db_instance.db_shard0.address}"
  value = "${aws_db_instance.db_addr.address}"
  // also: id of subnet group: aws_db_subnet_group.db_subnet_group.id
}


  
