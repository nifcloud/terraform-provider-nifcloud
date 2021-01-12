provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_dhcp_config" "basic" {
    static_mapping {
        static_mapping_ipaddress = "192.168.2.10"
        static_mapping_macaddress = "00:00:5e:00:53:FF"
        static_mapping_description = "static-mapping-memo-upd"
    }
}
