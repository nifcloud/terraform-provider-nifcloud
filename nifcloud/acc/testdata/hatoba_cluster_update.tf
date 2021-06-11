provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_hatoba_cluster" "basic" {
  name           = "%supd"
  description    = "memo-upd"
  firewall_group = nifcloud_hatoba_firewall_group.basic.name
  locations      = ["east-21"]

  addons_config {
    http_load_balancing {
      disabled = false
    }
  }

  network_config {
    network_id = "net-COMMON_PRIVATE"
  }

  node_pools {
    name          = "default"
    instance_type = "medium"
    node_count    = 3
  }

  node_pools {
    name          = "highspec"
    instance_type = "large"
    node_count    = 1
  }
}

resource "nifcloud_hatoba_firewall_group" "basic" {
  name = "%s"
}
