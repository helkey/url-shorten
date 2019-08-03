URL Shortener Design Doc
=============
Uniform Record Locator (URL) shorteners are used to access Internet resources, by providing a short URL to a resource that is easily typed and compactly stored.
Well-known URL shorteners include:

  * Bitly: the most popular and one of the oldest URL shorteners is used by Twitter for inserting links in tweets.
  By 2016 they had shortened 26 billion URLs
  * TinyURL. A simple shortener that requires no sign-up and allows users to customize the keyword
  * Goo.gl (DISCONTINUED): URL shortener written and retired by Google

Most URL shortener use is free, but Bitly projects [revenue in the range of $100M](https://www.cnbc.com/2016/05/26/web-link-shortening-company-bitly-eyeing-100m-revenues.html)
  by [providing Enterprise features](https://www.slant.co/versus/2591/22693/~bitly_vs_tinyurl).

## Key Features
Comparing to the leading ULR shortening service, this design has:

  * Higher security (12-character) standard links instead of 7 characters (Bitly standard links), 8 characters (TinyURL links),
    or 10 characters (Bitly Facebook links).
  * Additional (14-character) security needed for gray-listed sensitive domains (Box, Dropbox, Google Maps, ...)
  * Scalability designed into architecture: Cloud-based worker system design, with orchestration for automatic scaling
  * Database sharding information encoded into shortened URLS for additional scalability


## User Security
The use of URL shorteners can compromise security as the purpose of URL shorteners is to
[reduce entropy of URLs used to specify websites](https://freedom-to-tinker.com/2016/04/14/gone-in-six-characters-short-urls-considered-harmful-for-cloud-services/)

The address space of shortened URLs can be scanned by adversaries to find URLs that reveal confidential customer information.
Google has discontinued their URL shortening service (perhaps because of these security issues), but provide clear warnings
about risks of using the service even though they no longer provide it. The major URL shortening services which continue to operate
do not provide similar warnings about security issues.

![](figs/GoogleShortenerHighlighted.png "Google Security Warning")

As a result of reduced URL shortened address space, these addresses can be searched to find web sites containing:

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
that are using Bitly to provide shortened URLs.

![](figs/ScanningCost.png "URL Scanning Cost")

In contrast, this URL shortening project uses 12 characters for the standard baseline
security level, which is projected to cost ~$600M to scan in 2022.
Shortened URLs for sensitive domains use 14-character addresses, where scanning the entire URL space
is projected to cost ~$37B in 2022.

In addition, URL shortening services may use sequential codes for the shortened URLs, which further reduces security by
allowing recipients of a shorted URL to access compromised related URLs. Bitly appears to use a 6 character
URL shortening space for addresses shortened at a similar time. If someone finds a sensitive shortened Bitly URL,
they can scan all of the other URLs shorted around the same time for a [few hundred dollar](https://arxiv.org/pdf/1604.02734v1.pdf).


## URL Encoding
Counter-based, pseudo-random sequence for increased security

Scalability needs to be accommodated, and the solution needs to support in-service scaling to higher levels of traffic.

The length of shortened URLs needs to be long enough to provide unique results for every URL shortening request.
In this URL shortening architecture, shortened URLs will be constructed with characters a-z, A-Z, and 0-9,
for a total of 62 different characters (the same character set used by Bitly for shortening).

As an example, examine the address size to support a traffic load of 200 shortened requests/second,
over a period of 5 years. Seven character URLs would sufficient to meet this initial traffic estimate.
For comparison, standard shortened URLs are 7 characters for standard Bitly links, and 8 characters for TinyURL links.
For the reasons listed above, such short URLs should be wholy inaddequate for security purposes, so longer URLs are used.

### Grey-Listing Sensitive URLs
Sensitive URLs like Dropbox URLs or Maps URLs should not be as short as URLs suitable for public access.
In this application, URLs from these sensitive domains are gray listed for special processing, initially shortened to 12 characters
(for an address space of 3 x 10^21) rather than 10 character addresses for less sensitive URLs.

Scanning a 12-character address space should increase the cost for a full scan from $37k to $34B (in 2016 prices),
which would seem to be sufficiently expensive to make URL scanning unattractive compared to exploiting vulnerabilities
in competing URL shorting services which are less well protected.

### Encode Architecture
Load balancer

### Decode Architecture
Database sharding makes decoding significantly easier.

### Address Range Server / Database
Highly reliability
Zookeeper is a [highly reliable distributed datastore](https://aphyr.com/posts/291-call-me-maybe-zookeeper) that is suitable for this task.
Zookeeper uses majority quorums - using five notes, any two nodes could fail without degrading the system.
Zookeeper is also linearizable - all clients see the same ordering for updates occurring in the same order.

### URL Database
In a commercially successful URL shortener, the service has competition offering free services which puts some limit on customer value.
A URL shortener without a free tier probably could not compete successfully with services providing
a free tier. A free URL shortening service acts as advertising for enterprise customers,
as most users become familiar with the service when copy/pasting in shortened links provided by others.

The data storage requirements are large even for a moderately successful player in this space.
For the modest example goal of 200 URL shortening request/sec would result over a 5 year period in

    200 * 60sec * 60min * 24hr * 365day * 5year ~ 32 billion shortening requests

The URL database needs to be reliable, but cost is a dominant issue.
The volume of database writes and reads is high, so operational costs are high.
As a result, a managed database like AWS RDS might not be commercially feasible for this application.

The required URL mapping could be implemented in a distributed key-value database.
A promising distributed database might be the open source *etcd* database, written in Go,
which uses the Raft consensus algorithm. Other candidate [distributed key-value
stores](https://www.g2.com/categories/key-value-stores) include Aerospike, ArandoDB, BoltDB, CouchDB, Google Cloud Datastore, Hbase, and Redis.

### URL Database Sharding / Replicas
Database access to store the mapping from shortened URLs to URLs can be a
bottleneck for performance, limiting the scalability of popular web-based application.

Database sharding allows the generated data to be split across multiple databases,
reducing the load on each database. Here database sharding is implemented as a key
part of the software architecture. Database sharding, together with scalable cloud workers
to scale resources, should allow this URL shortener implementation to scale to 
levels of use similar to commercial competitors (Bit.ly and others)

[Database read replicas](https://aws.amazon.com/rds/details/read-replicas/) are useful for read-heavy applications.
A URL shortening application is an ideal case where database reads and writes from separate applications.


## Initial Implementation
A scalable URL shortening algorithm was implemented and deployed, using three AWS server instances for
the URL shortening server, URL expanding server, and internal address server. Two URL databases were deployed
to test database sharding, with an additional database used for the internal address server.

All of the server instances and databases in this demonstration were deployed to a default AWS public subnet,
which allowed easy external access to all of the deployed instances for testing.
***DRAWING**

### AWS Cloud Platform
Amazon Web Services (AWS) was chosen for the initial implementation, as AWS has over half the
cloud computing provider market share, and as a result the most software solutions. However,
Azure (Microsoft) and Google Cloud Compute would also have been reasonable choices.

### Amazon Machine Images
AWS uses Amazon Machine Images (AMIs) to customize cloud instances.
The Amazon drawing below shows the lifecycle of an AWS AMI.
AMIs can be created and stored, then registered when ready for instantiation.
Multiple identical cloud works can be generated from an AMI.
AMIs can be deregistered when no longer needed.

![](figs/ami_lifecycle.png "AMI life cycle")
[AMI lifecycle](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html]

AWS provides generic AMIs as a starting point for customization, such as the [Linux 2 AMI](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/finding-an-ami.html).
AMIs can be manually configured, for example by launching a generic AMI, customizing from the command line, then save updated the updated AIM with a new name.

AMIs can be tagged for identification, such as 'owner', 'development/production', or 'release number'.
Tags can help organize your AWS bill, for example for budgeting and accounting purposes.

AWS provides a [checklist](docs.aws.amazon.com/marketplace/latest/userguide/best-practices-for-building-your-amis.html) for AMIs, including:

    Linux-based AMIs that a valid SSH port is open (default is 22)
    ...
    
### Infrastructure as Code
Cloud computing resources can be configured through a GUI, but this process is error-prone,
and produces results that cannot always be repeated. It is far better to specify cloud resources
with code, which allows for version control, releases, version rollback, and many other features.

Containers have become a popular for server configuration, ensuring consistency between development and release cycles,
and between local testing and cloud-based deployment. Docker is a container solution often used to provide platform independence
and ease of managing resources.

The server software for this project consists of a single Go binary per instance. The Go language packages all dependencies
into a single executable, and seems like an adequate solution for this URL shortening application without adding
support for production containers.

### Packer
[Packer](https://www.packer.io/) is used for generating the AWS cloud instance images for this project.
Amazon AMIs could be manually configured, but should be configured as code for the reasons listed above
(e.g. documentation and repeatability).

For complex configurations, Packer interfaces with popular configuration management tools such as 
Chef and Ansible.

Packer (like Docker) allows building a single configuration management supporting multiple target platforms,
with support for cloud environments such as AWS, Azure, Google Cloud, as well as
desktop environments such as Linux, Windows, and Macintosh environments.
Packer also ntegrates with software platforms such as Docker, Kubernetes, and VMware.

### Terraform
Terraform is open-source cloud resource infrastructure deployment software.
It was chosen for this project based on its popularity and support of  major cloud platforms.

Terraform is programmed using functions directly mapped from cloud vendor modules,
and thus its provisioning code is not portable across cloud hardware vendors.
[Amazon CloudFormation](https://aws.amazon.com/cloudformation/) would also have been a reasonable choice
for configuring AWS-specific solution infracture, but Terraform does have a lot of features that
are portable across vendor platforms.

Terraform uses a declarative method for specifying deployments, which clearly documents
existing state of deployed infrastructure. Declarative provisioning is easier to operate
correctly and reliably than using a procedural approach, for example as provided by Chef and Ansible.

The open source Terraform version stores infrastructure configuration in a file on the users computer,
which is unsuitable for team use as  more than one team member might be making configuration changes at
the same time. In addition, the Terraform state file can contain sensitive encryption key information, and
so should not be put into a version control system. With small teams, Terraform can be used by putting the
infrastructure state file in shared, encrypted storage, and adding file locking to avoid simultaneous changes.

[Terraform Cloud](https://www.hashicorp.com/products/terraform/offerings) provides additional team functions for free,
including allowing remote infrastructure state access and locking during Terraform operations.
In addition, Terraform Cloud provides a user interface with history of changes to state, as well as information about who made
infrastructure state changes.

[Terraform Enterprise](https://www.hashicorp.com/products/terraform/offerings) is a paid product providing additional
features for larger teams, including operations features such as team management and a configuration designer tool,
and governance features such as audit logging.

### AWS Infrastructure

```terraform
```

### Address and URL Databases
For convenience, the initial project uses databases implemented using AWS RDS PostgreSql,
which is well supported by Terraform.

```terraform
```

## Next Generation Architecture

The next steps for this project are to move the databases to a private subnet,
set up a network address translation server for access from the private subnet,
implement load balancers with multiple URL shorten and expand servers,
configure autoscaling for these servers, and setting up continuous integration (CI) for development.

Further scaling work could include implementing a caching interface, setting up server groups in multiple geographic zones,
utilizing lower cost AWS spot instance, and investigating less expensive database storage options.


### Private Subnet
AWS example

```terraform
```

### NAT server
A network address translation (NAT) server can provide a number of functions...

```terraform
```

### Load Balancer
Kubernetes, AWS (Docker??)

Caddylightweight ingress service
  https://www.ardanlabs.com/blog/2019/07/caddy-partnership-light-code-labs.html
ELB (AWS)::
HAProxy
nginx

```terraform
```

### Bastion Host
The architecture used here has a public-facing load-balancer, which forwards traffic to the other worker nodes
  which are on a private network for improved security.

Management connections to the worker nodes is through a [bastion host](), which forwards (ssh) traffic to the worker nodes.
Load balancer as bastion host.

For improved security, a bastion host can be configured as a stand-alone server with a stripped-down operating system,
such as [AWS AppStream](https://aws.amazon.com/blogs/security/how-to-use-amazon-appstream-2-0-to-reduce-your-bastion-host-attack-surface/).

Bastion host
use ssh agent configured with agent forwarding from the client computer to avoid having to
  [store the private key on the bastion computer](https://aws.amazon.com/blogs/security/securely-connect-to-linux-instances-running-in-a-private-amazon-vpc/)




```terraform
```

### Autoscaling AWS Instances
AWS and other platforms have direct support for autoscaling features

```terraform
```

### Caching
Caching is another performance enhancing feature which is important (not included in current algorithm implementation).
A caching system (such as Redis or Memcache) will save responses to recent queries.
When requesting a shortened URL, caching will intercept frequent requests to provide shortened URLs for the same long URL,
which in turn minimizes wasted data storage due to multiple shortened versions of the same URL that would otherwise occur.
Caching also reduces load on the URL database when many users are requesting access to the same shortened URL,
by storing common resent requests in cache.

### AWS Spot Instances
AWS has reserved instances which can be held and operated as long as you like.
They also offer much less expensive spot instances, which are priced on an instantenous
'spot' model, and can preempted whenever someone bids a higher price for them.

When using spot instances to handle traffic overlow, reserved instances should still be
maintained at a baseline level to ensure meeting service level availability objectives.

In addition, instances need to handle the `shutdown -h` termination message promptly and
reliabably.

### Database alternatives

### Continuous Integration
CircleCI

CI/CD from GitHub to EC2 using AWS CodeBuild & CodeDeploy

## Infrastructure Orchestration
Terraform was used here for provisioning cloud resources due to its simplicity and ease of use.
Schedulers popular for large and complex systems include Fleet, Kubernetes, Marathon, and Mesos.

### Kubernetes
Kubernetes provides deployment, scaling, load balancing, and monitoring.
Kubernetes was developed at Google, and has become an extremely popular recently due to its power and flexibility.
In 2015, container survey found just 10 percent of respondents were using any container orchestration tool.
Two years, 71% of respondents were [using Kubernetes to manage their containers]().

Kubernetes is particularly well suited for a hybrid server use case, for example where some of the resources
are in an on-prem data center, and other resources are in the cloud.







