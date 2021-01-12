provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_dhcp_option" "basic" {
  default_router = "192.168.0.1"
  domain_name = "example.com"
  domain_name_servers = ["192.168.0.1", "192.168.0.2"]
  ntp_servers = ["192.168.0.1"]
  netbios_name_servers = ["192.168.0.1", "192.168.0.2"]
  netbios_node_type = "1"
  lease_time = "600"
}
