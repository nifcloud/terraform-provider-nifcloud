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

resource "nifcloud_dhcp_config" "example" {
    static_mapping {
        static_mapping_ipaddress = "192.168.1.10"
        static_mapping_macaddress = "00:00:5e:00:53:00"
        static_mapping_description = "static-mapping-memo"
    }
    ipaddress_pool {
        ipaddress_pool_start = "192.168.2.1"
        ipaddress_pool_stop = "192.168.2.100"
        ipaddress_pool_description = "ipaddress-pool-memo"
    }
}
