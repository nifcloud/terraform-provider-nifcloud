provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_db_parameter_group" "basic" {
    name       = "%s"
    family      = "mysql5.6"
    description = "tfacc-memo"

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
