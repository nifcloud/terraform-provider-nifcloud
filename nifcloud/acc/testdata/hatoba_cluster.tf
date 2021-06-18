provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_hatoba_cluster" "basic" {
  name           = "%s"
  description    = "memo"
  firewall_group = nifcloud_hatoba_firewall_group.basic.name
  locations      = ["east-21"]

  addons_config {
    http_load_balancing {
      disabled = true
    }
  }

  network_config {
    network_id = "net-COMMON_PRIVATE"
  }

  node_pools {
    name          = "default"
    instance_type = "medium"
    node_count    = 1
  }

  node_pools {
    name          = "lowspec"
    instance_type = "e-medium"
    node_count    = 1
  }
}

resource "nifcloud_hatoba_firewall_group" "basic" {
  name = "%s"
}
