// Addr database

// Use `google_sql_user` to define user host, password
resource "google_sql_database_instance" "url" {
  name = "db-url"
  database_version = "POSTGRES_9_6"
  region = "us-west1"
  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "url" {
  name     = "postgres"
  instance = "${google_sql_database_instance.url.name}"
  password = "${var.db_password}"   // Password is stored in TF state file. Encrypt state file,
                                    //  or modify afterward 
}
