variable "database_url" {
  type    = string
  default = getenv("DATABASE_URL")
}

variable "prod_database_url" {
  type    = string
  default = getenv("PROD_DATABASE_URL")
}

env "local" {
  url     = var.database_url
  dev     = "docker://postgres/18/dev?search_path=public"
  schemas = ["public"]

  migration {
    dir = "file://migrations"
  }
}

env "prod" {
  url     = "${var.prod_database_url}&search_path=public"
  dev     = "docker://postgres/18/dev?search_path=public"
  schemas = ["public"]

  migration {
    dir = "file://migrations"
  }
}
