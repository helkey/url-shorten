# Specify the provider and access details
provider "google" {
  // credential file stored in 
  credentials = "${file("urlshorten-2505-4fb8703bfc74.json")}"
  project     = "urlshorten-2505"
  region = "${var.gcp_region}"
}


