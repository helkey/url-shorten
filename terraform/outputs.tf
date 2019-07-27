// $terraform output shard_0 shard_1... > "aws.ip.txt"

//output "elb_address" {
// value = "${aws_elb.web.dns_name}"
//}

output "aws_key_name" {
  value="${aws_key_pair.auth.key_name}"
}

output "aws_public_key" {
  value="${aws_key_pair.auth.public_key}"
}
  
output "addr_ip" {
  value="${aws_instance.addr_server.public_ip}"
}

// output "addr_shard0" {
output "addr_dbAddr" {
  // value = "${aws_db_instance.db_shard0.address}"
  // also: id of subnet group: aws_db_subnet_group.db_subnet_group.id
  value = "${aws_db_instance.db_addr.address}"
}

