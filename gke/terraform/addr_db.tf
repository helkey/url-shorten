// Addr database

resource "google_sql_database_instance" "addr" {
  name = "master-instance"
  database_version = "POSTGRES_9_6"
  region = "us-west1"
  settings {
    tier = "db-f1-micro"
  }
}


