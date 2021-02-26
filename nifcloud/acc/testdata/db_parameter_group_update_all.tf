provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_db_parameter_group" "basic" {
    name       = "%supd"
    family      = "mysql5.7"
    description = "tfacc-memo-upd"

    parameter {
        name  = "character_set_server"
        value = "ascii"
    }

    parameter {
        name  = "character_set_results"
        value = "ascii"
    }
}
