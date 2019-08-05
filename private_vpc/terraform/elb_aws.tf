// https://www.terraform.io/docs/providers/aws/r/elb.html

# Create a new load balancer
resource "aws_elb" "url" {
  name               = "foobar-terraform-elb"
  availability_zones =["${var.aws_region}"]

  access_logs {
    // bucket        = "foo"
    bucket_prefix = "url"
    interval      = 60
  }

  listener {
    instance_port     = "${var.port_internal}"
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  listener {
    instance_port      = "${var.port_internal}"
    instance_protocol  = "http"
    lb_port            = 443
    lb_protocol        = "https"
    ssl_certificate_id = "arn:aws:iam::**TBD**:server-certificate/certName"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 4
    target              = "HTTP:" + "${var.port_internal}"
    interval            = 25
  }

  instances                   = ["${aws_instance.url.id}"]
  cross_zone_load_balancing   = true
  idle_timeout                = 500
  connection_draining         = true
  connection_draining_timeout = 500

  tags = {
    Name = "url_elb"
  }
}
