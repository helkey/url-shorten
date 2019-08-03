URL Shortener Design Doc
=============
Uniform Record Locator (URL) shorteners are used to access Internet resources, by providing a short URL to a resource that is easily typed and compactly stored.
Well-known URL shorteners include:

  * Bitly: the most popular and one of the oldest URL shorteners is used by Twitter for inserting links in tweets.
  By 2016 they had shortened 26 billion URLs
  * TinyURL. A simple shortener that requires no sign-up and allows users to customize the keyword
  * Goo.gl (DISCONTINUED): URL shortener written and retired by Google


## Key Features
Comparing to the leading ULR shortening service, this design has:

  * Higher security (12-character) standard links instead of 7 characters (Bitly standard links), 8 characters (TinyURL links),
    or 10 characters (Bitly Facebook links).
  * Additional (14-character) security needed for gray-listed sensitive domains (BoxDropbox, Google Maps, ...)
  * Scalability designed into architecture: Cloud-based worker system design, with orchestration for automatic scaling
  * Database sharding information encoded into shortened URLS for additional scalability


## User Security
The use of URL shorteners can compromise security as the purpose of URL shorteners is to reduce entropy of URLs used to specify websites [Shmatikov][Shmatikov-blog].
The address space of shortened URLs can be scanned by adversaries to find URLs that reveal confidential customer information.
Google has discontinued their URL shortening service (perhaps because of these security issues), but provide clear warnings
about risks of using the service even though they no longer provide it. The major URL shortening services which continue to operate
do not provide similar warnings about security issues.

![](figs/GoogleShortenerHighlighted.png "Google Security Warning")

As a result of reduced URL shortened address space, these addresses can be searched to find web sites containing

  * Cloud storage URLs for documents such as Box, Dropbox, GoogleDrive, and OneDrive documents.
      This is a <i>huge</i> security issue. For instance OneDrive links not only let adversaries edit
      the document, they can also use this link to gain access to other files [Shmatikov].
  * Map trip description URLs which may include the users identifiable home address linked to destinations.
      By starting from an address and mapping all endpoints from multiple URLs, one can create a
      personal connection graph by determining who visited whom [Shmatikov].

URL shorteners should provide shortened versions that are long enough to make adversarial scanning unattractive,
  limit the scanning of large numbers of potential URLs (by CAPTCHAS and IP blocking),
  and avoid generation of sequential URL addresses.
  
The adversarial cost of scanning the standard 7-bit Bit.ly address space was $37k in 2016 [Shmatikov]. 
The cost of Internet transit dropped 36% per year from 2010-2015 [BandwidthPriceTrends].
Using these two data points, we can project that by 2022 it will be possible to scan
all of a 10-character URL space for around $10M, so even the highest security level
that Bitly offers is not good enough for securing the large number of sensitive URLs
that are using Bitly to provide shortened URLs.

![](figs/ScanningCost.png "URL Scanning Cost")

In contrast, this URL shortening project uses 12 characters for the standard baseline
security level, which is projected to cost ~$600M to scan in 2022.
Shortened URLs for sensitive domains use 14 character address, where scanning the entire URL space
is projected to cost ~$37B in 2022.

In addition, URL shortening services may use sequential codes for the shortened URLs, which further reduces security by
allowing receipients of a shorted URL to access compromised related URLs. Bitly appears to use a 6 character
URL shortening space for addresses shortened at a similar time. If someone finds a sensitive shortened Bitly URL,
they can scan all of the other URLs shorted around the same time for a few hundred dollars [Shmatikov].
the majority of the users.

## Design
The length of shortened URLs needs to be long enough to provide unique results for every URL shortening request.
In this URL shortening architecture, shortened URLs will be constructed with characters a-z, A-Z, and 0-9,
for a total of 62 different characters (the same character set used by Bitly for shortening).

As an example, examine the address size to support a  traffic load of 200 shortened requests/second,
over a period of 5 years. Seven character URLs would sufficient to meet this initial traffic estimate.
For comparison, standard shortened URLs are 7 characters for standard Bitly links, and 8 characters for TinyURL links.
For the reasons listed above, such short URLs sould be wholy inaddequate for security purposes, so longer URLs are used.


### Architecture
Counter-based, pseudo-random sequence for increased security

Scalability needs to be accommodated, and the solution needs to support in-service scaling to higher levels of traffic.

### Caching
Caching is another performance enhancing feature which is important (not included in current algorithm implementation).
A caching system (such as Redis or Memcache) will save responses to recent queries.
When requesting a shortened URL, caching will intercept frequent requests to provide shortened URLs for the same long URL, which in
turn minimizes wasted data storage due to multiple shortened versions of the same URL that would otherwise occur.
Caching also reduces load on the URL database when many users are requesting access to the same shortened URL, by storing common
resent requests in cache.

### Load Balancer
Kubernetes, AWS (Docker??)
Caddylightweight ingress service
  https://www.ardanlabs.com/blog/2019/07/caddy-partnership-light-code-labs.html


### Encode Architecture
Load balancer

### Decode Architecture
Database sharding makes decoding significantly easier.


### Address Range Database
Highly reliability
Zookeeper is a [highly reliable distributed datastore](https://aphyr.com/posts/291-call-me-maybe-zookeeper) that is suitable for this task.
Zookeeper uses majority quorums - using five notes, any two nodes could fail without degrading the system.
Zookeeper is also linearizable - all clients see the same ordering for updates occuring in the sam.


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
A promising distributed database might be the open-source *etcd* database, written in Go,
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

### System Capacity Scalability
The system design should be scalable to handle more success than budgeted.
A successful URL shortening project could result in far more than 200 shortening requests per second.
Successful operation of the project would result in more than 5 years of operation.
As a result, in-service expansion of capability is highly desirable.

In the case of higher than expected demand, the URL encoding function should be upgraded to add additional
characters (from 7 to 8-10 characters). In the case of system capacity expansion,
the URL decoding software will chose which generation of URL encoding schema to apply.
Each time system capacity is increased by increasing number of characters,
there also will be an option to increase the number of database shards.


### Grey-Listing Sensitive URLs
Sensitive URLs like Dropbox URLs or Maps URLs should not be shortened like URLs suitable for public access.
Here URLs from these sensitive domains are gray-listed for special processing, initially shortened to 12 characters
(for an address space of 3x10^21) rather than the initial 7 character shortened addresses.
Scanning a 12-character address space should increase the full scanning cost from $37k to $34B (in 2016 prices),
which would seem to be sufficiently high to make URL scanning unattractive compared to exploiting vulnerabilities
in other URL shorting services.




## AWS Implementation
Load balancer

## Load Balancer
ELB (AWS)::
HAProxy
nginx

### Database
For convenience, the initial demonstration uses databases implementated using AWS RDS PostgreSql,
which is well supported by Terraform.




## Network Security
The architecture used here has a public-facing load-balancer, which forwards traffic to the other worker nodes
  which are on a private network for improved security.

Management connections to the worker nodes is through a [bastion host](), which forwards (ssh) traffic to the worker nodes.
Load balancer as bastion host.

For improved security, a bastion host can be configured as a stand-alone server with a stripped-down operating system,
such as [AWS AppStream](https://aws.amazon.com/blogs/security/how-to-use-amazon-appstream-2-0-to-reduce-your-bastion-host-attack-surface/).


## Cloud-Based Application Provisioning

### Terraform
Terraform was used to provision cloud worker resourse for this URL shortener.
Terraform was chosen as a provisioning solution based on being an open source, cloud agnostic provisioning tool.
Terraform uses a declarative method for specifying deployments, which clearly documents existing state of deployed infrastructure.
Declarative provisioning is easier to operate correctly and reliably than using a procedural approach, for example as provided by Chef and Ansible.

Terraform can also be used to allocate database resources.

Bastion host
use ssh agent configured with agent forwarding from the client computer to avoid having to
  [store the private key on the bastion computer](https://aws.amazon.com/blogs/security/securely-connect-to-linux-instances-running-in-a-private-amazon-vpc/)


## Configuration Management
Containers have become popular for ensuring consistency between development and release cycles,
and local and cloud-based deployment. Docker is often used for its benefits of container management,
including platform independence and ease of managing resources. The solution here runs a single Go binary
on each instance. Go packages all dependances into a single executable and so a simpler solution seems appropriate.

### Packer
[Packer]() is used here for generating cloud instance images used by Terraform.
Packer ensures that the images you test on are the same as deployed for production.

Packer (like Docker) allows building a single configuration that supports
multiple target platforms, and supports cloud environments such as AWS, Azure, Google Cloud, as well as
desktop environments such as Linux, Windows, and Macintosh environments, and platforms such as Docker, Kubernetes, and VMware.

For complex configuration management, Packer interfaces with popular configuration management tools such as 
Chef and Ansible. Go makes configuration management simple by packaging up the executable into a single binary.
Here simple Packer configuration scripts were used for configuration management.


## Cloud Hosting
AWS

### Amazon Machine Images
AWS uses Amazon Machine Images (AMIs) to customize cloud instances.
The drawing below shows the lifecycle of an AWS AMI.
AMIs can be created and stored, then registered when ready for instantiaion.
Multiple identical cloud works can be generated from an AMI.
The AMI should be deregistered when no longer needed.

![](figs/ami_lifecycle.png "AMI life cycle")
[AMI lifecycle](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html]

AMIs can be tagged for identification, such as 'owner', 'development/production', or 'release number'.
Tags to organize your AWS bill, for example for budgeting and accounting purposes.

AWS provides generic AMIs as a starting point for customization, such as the [Linux 2 AMI](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/finding-an-ami.html).
AMIs can be manually configured, for example by launching a generic AMI, customizing from the command line, then save updated the updated AIM with a new name.

AWS provides a [checklist](docs.aws.amazon.com/marketplace/latest/userguide/best-practices-for-building-your-amis.html) for AMIs, including:

    Linux-based AMIs that a valid SSH port is open (default is 22)
    ...
    

## Container orchestration

### Kubernetes
Kubernetes provides deployment, scaling, load balancing, and monitoring.
Kubernetes was developed at Google, and has become an extremely popular recently due to power and flexibilty.
In 2015, container survey found just 10 percent of respondents were using any container orchestration tool.
Two years, 71% of respondents were using Kubernetes to manage their containers.

Kubernetes is particularly well suited for a hybrid use case, where some of the resources are in an on-prem data center,
and other resources are in the cloud. Kubernetes is complex, and this application could be scaled with a simpler solution.


### Terraform
Terraform was used here for provisioning cloud resources due to its simplicity and ease of use.
It is well designed, and easy to use.

Terrafrom can also be used for orchestrating resources based network traffic.
Other schedulers popular for large and complex systems include Fleet, Marathon, Mesos, and Kubernetes.

## Continous Integration
CI/CD from GitHub to EC2 using AWS CodeBuild & CodeDeploy


## Testing



[Shmatikov]: [Gone in Six Characters: Short URLs Considered Harmful for Cloud Services](https://arxiv.org/pdf/1604.02734v1.pdf)

[Shmatikov-blog]: [(blog post) Gone In Six Characters: Short URLs Considered Harmful...](https://freedom-to-tinker.com/2016/04/14/gone-in-six-characters-short-urls-considered-harmful-for-cloud-services/)

[BandwidthPriceTrends]: [Internet transit pricing](http://drpeering.net/white-papers/Internet-Transit-Pricing-Historical-And-Projected.php)

[FeatureComparison]: [Bit.ly vs TinyURL](https://www.slant.co/versus/2591/22693/~bitly_vs_tinyurl)

[BitlyRevenue] https://www.cnbc.com/2016/05/26/web-link-shortening-company-bitly-eyeing-100m-revenues.html]