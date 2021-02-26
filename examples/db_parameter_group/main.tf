terraform {
  required_providers {
    nifcloud = {
      source = "nifcloud/nifcloud"
    }
  }
}

provider "nifcloud" {
  region = "jp-east-1"
}

resource "nifcloud_db_parameter_group" "default" {
    name        = "example"
    family      = "mysql5.7"
    description = "memo"

    parameter {
        name  = "character_set_server"
        value = "utf8"
    }

    parameter {
        name  = "character_set_client"
        value = "utf8"
    }

    parameter {
        name         = "character_set_results"
        value        = "utf8"
        apply_method = "pending-reboot"
    }
}
