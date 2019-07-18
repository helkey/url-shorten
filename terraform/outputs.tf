// $terraform output shard_0 shard_1... > "aws.ip.txt"

//output "elb_address" {
// value = "${aws_elb.web.dns_name}"
//}

output "addr_shard_0" {
  value = "${aws_db_instance.shard0.address}"
  //     = "${aws_db_instance.shard0.id}" # id
  //     = "${aws_db_instance.shard0.endpoint}" # endpoint w/port
  // id of subnet group: aws_db_subnet_group.db_subnet_group.id
}

output "id_shard_0" {
  value = "${aws_db_instance.shard0.id}"
}

output "endpoint_shard_0" {
  value = "${aws_db_instance.shard0.endpoint}"
}


  
