# url-shorten
# Markdown: www.cheatography.com/lucbpz/cheat-sheets/the-ultimate-markdown/

URL Shortener Design Doc
=============
URL shorteners are used to provide a short URL for accessing web pages that is easily typed and compactly stored.
Well-known examples include:
* Bitly: the most popular and one of the oldest URL shorteners is used by Twitter for inserting links in tweets.
  By 2016 they had shortened 26 billion URLs
* Goo.gl: URL shortener from Google
* TinyURL. A simple shortener that requires no sign-up and allows users to customize the keyword

Security
------­-
The use of URL shorteners can compromise security as reducing entropy of URLs used to specify websites is the purpose of the function [Shmatikov][Shmatikov-blog].
As a result, the URL shortened address space can be searched to find web sites containg
* Cloud storage URLs for documents such as Box, Dropbox, GoogleDrive, and OneDrive documents.
  This can a <b>huge</b> security issue. For instance OneDrive links not only let adversaries edit
  the document, they can also use this link to gain access to other files [Shmatikov].
* Map trip description URLs which may include the users identifiable home address linked to destinations.
  By starting from an address and mapping all enpoints, one can create a map of who visited whom [Shmatikov].

URL shorteners may use sequential codes for the shortened URLs, which further reduces security by
allowing receipts of a shorted URL to accss compromised related URLs.

URL shorteners should provide shortened URLs that are long enough to make advesarial scanning unattractive,
  limit the scanning of large numbers of potential URLs (by CAPTCHAS and IP blocking),
  and avoid generation of sequential URL addresses.
  
Users of URL shorteners should avoid using public URL shorteners for sensitive websites that could be compromised by adversaries,
  as well as taking care when clicking on shortened links which may take them to malicous websites.

Minimum Viable Product Features
------­-----
This implementation is designed as a minimum demonstration of a URL shortening application,
with an emphasis on scalability.

Requirements
------­-----
The length of shortened URLs needs to be long enough to provide unique results for every URL shortening request.
Here we wll assume fixed length short URLs, although a more general architecture would increase length of the URLs
as the URL space filled up.

For this design, a traffic load of 1000 shortened request/second, over a period of 5 years. In addition,
a 100x increase in the URL address space is allocated in order to scanning of the shortened URL
address space unattractive compared to scanning other URL shortening services. This is similar to the
principal that you can't make your property impervious the theft, but you can make it a less attractive
than your neighbor's propert.

Design
------­-----
This application designed using Docker for its benefits of container management, including platform independance and ease of managing resources,
and Kubernetes for load balancing and resource orchestration.

## Encode Architecture


## Decode Architecture


## Database Sharding
Database access to store the mapping from shortened URLs to URLs can be a
bottlenect for performance, limiting the scalability of popular web-based application.

Database sharding allows the generated data to be split across multiple databases,
reducing the load on each database. Here database sharding is implemented as a key
part of the software architecture. Database sharding, together with using Kubernetes
to scale resources, should allow this URL shortener implementation to scale to high
levels of use.

## Caching
Caching is another performance enhancing feature which is important (not included in current algorithm implementation).
A caching system (such as Redis or Memcache) will save responses to recent queries.
When requesting a shortened URL, caching will intercept frequent requests to provide shortened URLs for the same long URL, which in
turn minimizes wasted data storage due to multiple shortened versions of the same URL that would otherwise occur.
Caching also reduces load on the URL database when many users are requesting access to the same shortened URL, by storing common
resent requests in cache.

## URL Encoding
In this URL shortening architecture, shortened URLs will be constructed with characters a-z, A-Z, and 0-9, for a total of 62 different characters
(the same character set used by Bit.ly for shortening).

The system goal of 1000 URL shortening request/sec means over a 5 year period the expected number of requests is
200 * 60sec * 60min * 24hr * 365day * 5year ~ 32 billion requests
Assuming an overcapcity ratio of 20x to make scanning the URL space less attractive,
the number of URLs available ~ 3.2 trillion, which can be encoded using 7 characters (taken from the 62 character set).
200*60*60*24*365*5 / 1e9
100*200*60*60*24*365*5 / 1e12
log(100*200*60*60*24*365*5)/log(62)

## System Capacity Scalability
The system design should be scalabile to handle more success than budgeted.
A successful URL shortening project could result in far more than 200 shortening requests per second.
Successful operation of the project would result in more than 5 years of operation.
As a result, in-service expansion of capability is highly desirable.

In the case of higher than expected demand, the URL encoding function should be upgraded to add additional
characters (from 7 to 8-10 characters). In the case of system capacity expansion,
the URL decoding software will chose which generation of URL encoding schema to apply.
Each time system capacity is increased by increasing number of characters,
there also will be an option to increase the number of database shards.

## Grey-Listing Sensitive URLs
The cost of scanning the 7-bit Bit.ly address space was $37k in 2016, while the Bit.ly address space was about 37% used [Shmatikov].
Even by overallocating the URL address space by 100x, this is inexpensive enough that URL shortening systems should
assume that their full address space will be scanned looking for senstive information.

Sensitive URLs like Dropbox URLs or Maps URLs should not be shortened like URLs suitable for public access.
Here URLs from these sensitive domains are grey-listed for special processing, initially shortened to 12 characters
(for an address space of 3x10^21) rather than the initial 7 character shortened addresses.
Scanning a 12-character address space should increase the full scanning cost from $37k to $34B (in 2016 costs),
which would seem to be suffiently high to make URL scanning unattractive compared to exploiting vulnerabilities
in other URL shorting services.

Performance Testing Results
------­-----

References
------­-----
[Shmatikov] Gone in Six Characters: Short URLs Considered Harmful for Cloud Services,
            https://arxiv.org/pdf/1604.02734v1.pdf

[Shmatikov-blog] Gone In Six Characters: Short URLs Considered Harmful for Cloud Services,
      https://freedom-to-tinker.com/2016/04/14/gone-in-six-characters-short-urls-considered-harmful-for-cloud-services/

