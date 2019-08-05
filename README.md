URL Shortener Design Doc
=============
Uniform Record Locator (URL) shorteners are used to access Internet resources, by providing a short URL to a resource that is easily typed and compactly stored.
Well-known URL shorteners include:

  * Bitly: the most popular and one of the oldest URL shorteners is used by Twitter for inserting links in tweets.
  By 2016 they had shortened 26 billion URLs
  * TinyURL. A simple shortener that requires no sign-up and allows users to customize the keyword
  * Goo.gl: URL shortener (DISCONTINUED SERVICE) written and
     [shut down](https://developers.googleblog.com/2018/03/transitioning-google-url-shortener.html) by Google

Most URL shortener use is free, but projections for Bitly revenue is in the
  [range of $100M](https://www.cnbc.com/2016/05/26/web-link-shortening-company-bitly-eyeing-100m-revenues.html),
  achieved by a freemium model with [paid Enterprise features](https://www.slant.co/versus/2591/22693/~bitly_vs_tinyurl).


## Key Features
In contrast to the leading ULR shortening service, the features of this design include:

  * Higher security (12-character) standard links instead of 7 characters (Bitly standard links), 8 characters (TinyURL links),
    or 10 characters (Bitly Facebook links).
  * Additional (14-character) security needed for gray-listed sensitive domains (Box, Dropbox, Google Maps, ...)
  * Scalability designed into architecture: Cloud-based worker system design, with orchestration for automatic scaling
  * Database sharding information encoded into shortened URLS for additional scalability


## User Security
The use of URL shorteners can compromise security as the purpose of URL shorteners is to
[reduce entropy of URLs](https://freedom-to-tinker.com/2016/04/14/gone-in-six-characters-short-urls-considered-harmful-for-cloud-services/)
used to specify websites.

The address space of shortened URLs can be scanned by adversaries to find URLs that reveal confidential customer information.
Perhaps because of these security issues, Google has discontinued their URL shortening service, but maintain service
to expand previously shortened URLs and provide clear warnings about risks of using the service even though they no longer provide it.
The other major URL shortening services which continue to operate do not provide warnings about security issues
in using URL shorteners.

![](figs/GoogleShortenerHighlighted.png "Google Security Warning")

As a result of reduced URL shortened address space, possible URL shortened addresses can be scanned to find web sites containing:

  * Cloud storage URLs for documents such as Box, Dropbox, GoogleDrive, and OneDrive documents.
      This is a <i>huge</i> security issue. For instance, OneDrive links not only let adversaries edit
      the document, they can also use this link to [gain access to other files](https://arxiv.org/pdf/1604.02734v1.pdf).
  * Map trip description URLs which may include the users identifiable home address linked to destinations.
      By starting from an address and mapping all endpoints from multiple URLs, one can create a
      personal connection graph by [determining who visited whom](https://arxiv.org/pdf/1604.02734v1.pdf).

URL shorteners should provide shortened versions that are long enough to make adversarial scanning unattractive,
  limit the scanning of large numbers of potential URLs (by CAPTCHAS and IP blocking),
  and avoid generation of sequential URL addresses.
  
The cost of adversarial [scanning the standard 7-bit Bit.ly address space](https://arxiv.org/pdf/1604.02734v1.pdf)
was $37k in 2016. The cost of Internet transit [dropped 36% per year from 2010-2015](http://drpeering.net/white-papers/Internet-Transit-Pricing-Historical-And-Projected.php)

Using these two data points, we can project that by 2022 it will be possible to scan
all of a 10-character URL space for around $10M, so even the highest security level
that Bitly offers is not good enough for securing the large number of sensitive URLs
that are using Bitly to provide URL shortening.

![](figs/ScanningCost.png "URL Scanning Cost")

In contrast, this URL shortening project uses 12 characters for the standard baseline
security level, which is projected to cost ~$600M for a full scan in 2022.
In addition, this project provides shortened URLs for sensitive domains use 14-character addresses,
where scanning the entire URL space is projected to cost ~$37B in 2022.

Another security vulnerability is that URL shortening services may use sequential codes for the shortened URLs,
which further reduces security by allowing recipients of a shorted URL to access compromised related URLs.
Bitly appears to use a 6-character URL shortening space for addresses shortened at a similar time.
If someone finds a sensitive shortened Bitly URL, they can scan all of the other URLs shortened around the same time for a
  [few hundred dollar](https://arxiv.org/pdf/1604.02734v1.pdf).


## URL Encoding
The length of shortened URLs needs to be long enough to provide unique results for every URL shortening request.
In this URL shortening architecture, shortened URLs will be constructed with characters a-z, A-Z, and 0-9,
for a total of 62 different characters (the same character set used by Bitly for short URLs).
As mentioned about, standard URL shortening provides 12-character URLs.

### Grey-Listing Sensitive URLs
Sensitive URLs like Dropbox URLs or Maps URLs should not be as short as URLs suitable for public access.
In this application, URLs from these sensitive domains are gray listed for special processing, initially shortened to 12 characters
(for an address space of 3 x 10^21) rather than 10-character addresses for less sensitive URLs.

Scanning a 12-character address space should increase the cost for a full scan from $37k to $34B (in 2016 prices),
which would seem to be sufficiently expensive to make URL scanning unattractive compared to exploiting vulnerabilities
in competing URL shorting services which are less well protected.

### Database Shard Encoding
Database sharding, where separate databases are used to encode different data, makes scaling of distributed databases more efficient.
In this project, a database shard is assigned to each shortened URL, allowing the expanded URL to be recovered
by querying a smaller database than the size that would be required without sharding.

### Address Range Server / Database
The initial implementation here uses shortening servers which provide URL shortening, and an address range server to provide each shortening
server with unique sets of addresses. Shortened URLs within this address range are served in random order.

As mentioned before, using address ranges is a potential security issue, allowing someone with a shortened URL to one of your resources
to more easily find other related resources shortened at a similar time. This security issue is mitigated here as different shortener
servers have different address ranges, and subsequent shortening requests will likely have completely unrelated addresses.
However, further work would be needed to deploy a commercial URL shortening system without any detectable correlation between addresses.

In order to allow distributed cloud instances to assign unique shortened URLs, a server is used to allocate encoded address
ranges to each instance. This address range server needs to be highly reliable in order to avoid assigning the same shortened URL codes
to multiple long URLs. Here a centralized server is used to generate small address ranges and assign them to 

A highly reliable distributed datastore such as [Zookeeper](https://aphyr.com/posts/291-call-me-maybe-zookeeper)
would be a better choice for this address range server task. Zookeeper uses majority quorums - using five notes,
any two nodes could fail without degrading the system. Zookeeper is also linearizable - all clients see the same ordering
for updates occurring in the same order.

### URL Database Sharding / Replicas
Database access to store the mapping from shortened URLs to full URLs can be a
bottleneck for performance, limiting the scalability of popular web-based application.

Database sharding allows the generated data to be split across multiple databases,
reducing the traffic load on each database. Here database sharding is implemented as a key
part of the software architecture. Database sharding, together with scalable cloud workers
to scale resources, should allow this URL shortener project to scale to levels of use
similar to commercial competitors (Bit.ly and others)

[Database read replicas](https://aws.amazon.com/rds/details/read-replicas/) are useful for read-heavy applications
such as this. A URL shortening application is an ideal case for this technology, where database reads and writes
come from separate applications.


## Initial Implementation
A scalable URL shortening algorithm was implemented and deployed, using three AWS server instances for
the URL shortening server, URL expanding server, and internal address server. Two URL databases were deployed
to test database sharding, with an additional database used for the internal address server.

All of the server instances and databases in this demonstration were deployed to a default AWS public subnet,
which allowed easy external access to all of the deployed instances for testing.

Go channels are used as buffers from the address server to the shortener servers,
in order account for data transmission errors which are more likely to occur in network applications.





### AWS Cloud Platform
Amazon Web Services (AWS) was chosen for the initial implementation, as Amazon has over half the
cloud computing provider market share, and as a result the most mature software solutions.
However, there are many other cloud computing choices. Azure (Microsoft) and Google Cloud Compute
are the next most popular cloud computing providers.

### Amazon Machine Images
AWS uses Amazon Machine Images (AMIs) to customize and manage cloud instances.
The Amazon drawing below shows the lifecycle of an AWS AMI.
AMIs can be created and stored, then registered when ready for instantiation.
Multiple identical cloud instances can be generated from an AMI.
AMIs can be deregistered when no longer needed to free up storage.

![](figs/ami_lifecycle.png "AMI life cycle")
[AMI lifecycle](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html]

AWS provides generic AMIs as a starting point for customization, such as the
  [Linux 2 AMI](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/finding-an-ami.html).

AMIs can be tagged for identification, such as 'owner', 'development/production', or 'release number'.
Tags can help organize your AWS bill, for example for budgeting and accounting purposes.

AWS provides a [checklist](docs.aws.amazon.com/marketplace/latest/userguide/best-practices-for-building-your-amis.html) for AMIs, including:

    Linux-based AMIs that a valid SSH port is open (default is 22)
    ...
    
### Infrastructure as Code
Cloud computing resources can be configured through a graphical user interface (GUI), but this process is error-prone,
and produces results that cannot always be repeated. It is far better to specify cloud resources with code,
which allows for version control, releases, version rollback, and many other features.

Containers have become a popular for server configuration, ensuring consistency between development and release cycles,
and between local testing and cloud-based deployment. Docker is a container solution often used to provide platform independence
and ease of managing resources.

The server software for this project consists of a single Go binary per instance. The Go language packages all dependencies
into a single executable, and seems like an adequate solution for this URL shortening application without adding
support for production containers.

### Packer
[Packer](https://www.packer.io/) is used in this project for generating the AWS cloud instance images.
AMIs can be configured manually, for example by launching a generic AMI, customizing from the command line,
  then save updated the updated AIM with a new name. However, it would be much better to configure AWS AMIs
as code for the reasons listed above (e.g. documentation and repeatability).

Packer (like Docker) allows building a single configuration management supporting multiple target platforms,
with support for cloud environments such as AWS, Azure, Google Cloud, as well as
desktop environments such as Linux, Windows, and Macintosh environments.
Packer also integrates with software platforms such as Docker, Kubernetes, and VMware.

For complex configurations, Packer interfaces with popular configuration management tools such as 
Chef and Ansible. In this example, each instance is running a single Go binary.
This server application is put in the `scripts/per-boot` folder, so that it runs
on startup, and after any reboot of the AWS instance.

```Packer
{
    ...
    "builders": [{
	"type": "amazon-ebs",
	"access_key": "{{user `aws_access_key`}}",
	"secret_key": "{{user `aws_secret_key`}}",
	"region": "us-west-1",
	"source_ami_filter": {
	    "filters": {
		"virtualization-type": "hvm",
		"name": "amzn2-ami-hvm-2.0.*-x86_64-gp2",
		"root-device-type": "ebs"
	    },
	    "owners": ["137112412989"],
	    "most_recent": true
	},
	"instance_type": "t2.micro",
	"ssh_username": "ec2-user",
	"ami_name": "ami-addr {{timestamp}}"
    }],
      "provisioners": [
    {
      "type": "file",
      "source": "../ReqAddr",
      "destination": "/tmp/ReqAddr"
    },
    {	  
      "type": "shell",
      "inline": [
        "sudo chmod 700 /tmp/ReqAddr",
        "sudo mv /tmp/ReqAddr /var/lib/cloud/scripts/per-boot/",
        "sleep 30",
        "sudo yum -y update"
      ]
    }]
}
```

T2.micro instances are used for development, as they are covered by the
[AWS free tier](https://aws.amazon.com/free/?all-free-tier.sort-by=item.additionalFields.SortRank&all-free-tier.sort-order=asc),
which provides 750 hours of small instance use per month for a year.


### Terraform
Terraform is open-source cloud resource infrastructure deployment software.
It was chosen for this project based on its popularity, and its support of major cloud platforms.

Terraform uses a declarative method for specifying deployments, which documents
existing state of deployed infrastructure, providing the advantages of infrastructure as code discussed earlier.
Declarative provisioning used by Terraform is easier to operate correctly and reliably
than by using a procedural programming approach, for example as provided by Chef and Ansible.

Terraform is programmed using functions directly mapped from the particular cloud vendor modules,
and as a result its provisioning code is not portable across cloud hardware vendors.
[Amazon CloudFormation](https://aws.amazon.com/cloudformation/) would also have been a reasonable choice
for configuring AWS-specific solution infrastruture, but Terraform does have a lot of features that
are portable across vendor platforms.

The open source Terraform version stores infrastructure configuration in a file on the user's computer,
which is unsuitable for team use as more than one team member might be making configuration changes at
the same time. In addition, the Terraform state file can contain sensitive encryption key information, and
so should not be put into a version control system. With small teams, Terraform can be used by putting the
infrastructure state file in shared, encrypted storage, and adding file locking to avoid simultaneous changes.

[Terraform Cloud](https://www.hashicorp.com/products/terraform/offerings) provides additional team functions for free,
including allowing remote infrastructure state access and locking during Terraform operations.
In addition, Terraform Cloud provides a user interface with history of changes to state, as well as information about who made
infrastructure state changes.

[Terraform Enterprise](https://www.hashicorp.com/products/terraform/offerings) is a paid product providing additional
features, including operations features such as team management and a configuration designer tool,
and governance features such as audit logging. Terraform Enterprise is strongly recommended for large teams
managing infrastructure with Terraform.

### AWS Infrastructure
Terraform uses a 'provider' to specify where/how to deploy the specified resources.
Here the region is specified in a variable 'aws_region'. Values for all of the variables
can be specified in a *.tfvars file, which allows easy support of multiple regions from
the same 
```terraform
// terraform/aws_provider.tf
provider "aws" {
  profile = "default"
  region = "${var.aws_region}"
  version = "~> 2.20"
}

resource "aws_key_pair" "auth" {
  key_name   = "${var.key_name}"
  public_key = "${file(var.public_key_path)}"
}

resource "aws_security_group" "instance" {
  name = "terraform-example-instance"

  // SSH access from anywhere
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  // other ports
  ...
```

### Address and URL Servers
Aws_instances are defined for each of the address, shortening, and expanding servicer instance.
As mentioned, small `t2.micro` AWS instances are used to minimize cost during development.

```terraform
resource "aws_instance" "addr_server" {
  ami = "${lookup(var.amis_addr, var.aws_region)}"
  instance_type = "t2.micro"
  key_name = "${var.key_name}"
  vpc_security_group_ids = [aws_security_group.instance.id]
```

### Address and URL Databases
For convenience, so far project uses databases implemented using AWS RDS PostgreSql,
which is supported by Terraform.

```terraform
resource "aws_db_instance" "db_shard0" {
  name                    = "db_shard0"
  allocated_storage       = 20 # GB
  engine                  = "postgres"
  instance_class          = "db.t2.micro"
  password = "${var.db_password}"
  port                    = 5433
  publicly_accessible     = true
  skip_final_snapshot     = true
  storage_type            = "gp2"
  username = "postgres"
  vpc_security_group_ids   = ["${aws_security_group.db.id}"]
```

## Next Generation Architecture
The next steps for this project are to move the databases from a public to a private subnet,
set up a network address translation server to allow internet access from the private subnet,
implement load balancers with multiple URL shorten and expand servers,
configure autoscaling for these servers to allow for variable traffic load,
and to set up continuous integration (CI) for development.

Further scaling work could include implementing a caching interface, setting up server groups in multiple geographic zones,
utilizing lower cost AWS spot instance, and investigating less expensive database storage options.


### Private Subnet
A private subnet is used to isolate functions that do not need public access,
in this case the address and URL databases, and potentially to isolate the address server as well.

The private network configuration includes defining the AWS subnet, route table, routing table association,
and one or more private network security groups.

Classless inter-domain routing (CIDR) blocks are used to define
  [IP address ranges](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_Subnets.html)
allowed for internal connections.

```terraform
resource "aws_subnet" "private" {
  vpc_id = "${aws_vpc.default.id}"
  cidr_block = "${var.private_cidr}"
  availability_zone = "us-west-1"
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
  ...
```

### NAT server
A network address translation (NAT) function is needed to provide access to the internet from
modules on the private subnet. 

Amazon provides pre-built AMIs for implementing a cloud instance configured as a NAT server.

```terraform
resource "aws_security_group" "nat" {
  name = "vpc_nat"
  description = "Allow traffic to pass from the private subnet to the internet"
  ingress {
    from_port = 22
    to_port = 22
    protocol = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
  ...
  
```

### Load Balancer
Two load balancer is needed to share traffic load to the short serves and expand servers.
AWS [elastic load balancing](https://aws.amazon.com/elasticloadbalancing/) (ELB) modules will be used here
for load balancing. There are a variety of other load balancing solutions that cold be using,
including setting up an Nginx server at the front-end.

```terraform
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
...
```

### Server autoscaling
A major advantage of cloud deployments is the ability to allocate resources on demand,
without paying for peak capability all of the time.

Amazon offers an autoscaling [autoscaling](https://www.terraform.io/docs/providers/aws/r/launch_configuration.html)
  group function which Terraform can configure.
  
```terraform
resource "aws_launch_configuration" "url" {
  image_id        = "ami-0c55b159cbfafe1f0"
  instance_type   = "t2.micro"
  security_groups = [aws_security_group.instance.id]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "auto" {
  name                 = "url_autoscaling_group"
  launch_configuration = "${aws_launch_configuration.url}"
  min_size             = 1
  max_size             = 8

  lifecycle {
    create_before_destroy = true
  }
}
...
```


### Bastion Host
The architecture used here has a public-facing load-balancer, which forwards traffic to compute instances for serving shortened or expanded URLs.
Network security can significantly improved by providing a single network entry point, known as a bastion host, for network control and monitoring.
Users can log in to the network using an SSH agent configured with agent forwarding from the client computer to avoid having to
  [store the private key on the bastion computer](https://aws.amazon.com/blogs/security/securely-connect-to-linux-instances-running-in-a-private-amazon-vpc/)

For maximum security, the bastion host can be configured as a stand-alone server providing only network access and traffic forwarding,
using a stripped-down operating system such as
  [AWS AppStream](https://aws.amazon.com/blogs/security/how-to-use-amazon-appstream-2-0-to-reduce-your-bastion-host-attack-surface/).

### Caching
Caching is another performance enhancing feature which is important feature to be implemented.
A caching system (such as Redis or Memcache) will save responses to recent queries.
When requesting a shortened URL, caching will intercept frequent requests to provide shortened URLs for the same long URL,
which in turn minimizes wasted data storage due to multiple shortened versions of the same URL that would otherwise occur.
Caching also reduces load on the URL database when many users are requesting access to the same shortened URL,
by storing common resent requests in cache.

### AWS Spot Instances
Cloud vendors such as AWS have reserved instances which can be held and operated indefinitely.
Vendors also offer much less expensive spot instances, which are priced on an instantaneous
'spot price' model, and can preempted whenever someone bids a higher price for them.

When using spot instances to handle traffic overflow, enough reserved instances should still be
maintained at a baseline level to ensure meeting service level availability objectives.
It addition, account limits on the maximum number of spot instances should be checked to verify
maximum capacity under heavy traffic.

Spot instances will be preempted frequently, so they need to handle a termination command (`shutdown -h`)
promptly and reliably. Switching from reserved instances to spot instances needs to be
[performed carefully](https://www.honeycomb.io/blog/treading-in-haunted-graveyards/) to avoid
introducing instance allocation issues.

### Continuous Integration
Continuous integration (CI) and continuous deliver (CD) are important for teams delivering high quality software.
AWS [CodeBuild](https://aws.amazon.com/codebuild/) and CodeDeploy provides CI/CD from GitHub code
to AWS EC2 server instances. [CircleCI](https://circleci.com/) is a very popular open-source CI/CD solution.

### Database alternatives
A commercially successful URL shortener service has existing competition offering free URL shortening,
which puts some limit on value can be extracted from the customer base.

A free URL shortening service is probably necessary, as it acts as advertising for enterprise customers,
and most users will become familiar with the service when copy/pasting in shortened links provided by others.
A URL shortener without a free tier probably could not compete successfully with successful services,
which to provide free use to most customer.

The data storage requirements are large even for a moderately successful player in this space.
For the modest example goal of 200 URL shortening request/sec would result over a 5-year period in

    200 * 60sec * 60min * 24hr * 365day * 5year ~ 32 billion shortening requests

The URL database needs to be reliable, but cost is a dominant issue.
The volume of database writes and reads is high, so operational costs are high.

Most SaaS applications have very good gross margins, and the cost of using a managed database
is a good tradeoff. However, a managed database like AWS RDS might not be commercially feasible
for this application, in which it should be difficult to maintain high gross margin.

The required URL mapping could be implemented in a distributed key-value database managed by internal staff.
Popular [distributed database](https://www.g2.com/categories/key-value-stores)
  candidates include Aerospike, ArandoDB, BoltDB, CouchDB, Google Cloud Datastore, Hbase, and Redis.
A promising distributed database written in Go is the open source *etcd* database, 
which uses the Raft consensus algorithm.


## Infrastructure Orchestration
Terraform was used here for provisioning cloud resources due to its simplicity and ease of use.
Schedulers popular for large and complex systems include Fleet, Kubernetes, Marathon, and Mesos.

### Kubernetes
Kubernetes provides deployment, scaling, load balancing, and monitoring.
Kubernetes was developed at Google, and has become an extremely popular recently due to its power and flexibility.
In 2015, container survey found just 10 percent of respondents were using any container orchestration tool.
Two years, 71% of respondents were
  [using Kubernetes to manage their containers](https://techcrunch.com/2017/12/18/as-kubernetes-surged-in-popularity-in-2017-it-created-a-vibrant-ecosystem/).

Kubernetes is particularly well suited for a hybrid server use case, for example the case
where some of the resources are in an on-prem data center, and other resources are in the cloud.







