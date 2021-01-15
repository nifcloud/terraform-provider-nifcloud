provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_instance" "basic" {
  instance_id       = "%s"
  availability_zone = "east-21"
  image_id          = "221"
  key_name          = nifcloud_key_pair.basic.key_name
  user_data         = <<EOT
#!/bin/bash

cat << EOS > /etc/netplan/99-netcfg.yaml
network:
  version: 2
  renderer: networkd
  ethernets:
      ens224:
          dhcp4: false
          addresses: [192.168.100.100/24]
          dhcp6: false
EOS
netplan apply
	EOT

  depends_on = [nifcloud_key_pair.basic, nifcloud_private_lan.basic]

  network_interface {
    network_name = nifcloud_private_lan.basic.private_lan_name
    ip_address   = "static"
  }

  network_interface {
    network_id = "net-COMMON_GLOBAL"
  }
}

resource "nifcloud_instance" "upd" {
  instance_id       = "%supd"
  availability_zone = "east-21"
  image_id          = "221"
  key_name          = nifcloud_key_pair.basic.key_name
  user_data         = <<EOT
#!/bin/bash

cat << EOS > /etc/netplan/99-netcfg.yaml
network:
  version: 2
  renderer: networkd
  ethernets:
      ens224:
          dhcp4: false
          addresses: [192.168.100.101/24]
          dhcp6: false
EOS
netplan apply
  EOT

  depends_on = [nifcloud_key_pair.basic, nifcloud_private_lan.basic]

  network_interface {
    network_name = nifcloud_private_lan.basic.private_lan_name
    ip_address   = "static"
  }

  network_interface {
    network_id = "net-COMMON_GLOBAL"
  }
}

resource "nifcloud_key_pair" "basic" {
  key_name   = "%s"
  public_key = "c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFEeFVVSmtIWFFvalVmeGphT3dQNVJmMjhOTVRFSjJFblBQdFk0b1NkZFBpRllnMWVDTGFNU08wV25nZVIrVk5sU215am1qU2xRWjBsc1BkcHZjWnY0KzZiMDlLUUZlT3NxakdjNE9Ga1o2MTZyTEI3UmdzblZnSXl3QmtIZ2lsMVQzbFRwRHVtYVk2TFFaRjRiaVpTNkNyaFdYeVhiSjFUVmYyZ0hIYXZPdi9WSS9ITjhIejlnSDg5Q0xWRVFOWFVQbXdjbC83ZE4yMXE4QnhNVkpGNW1sSW1RcGxwTjFKQVRwdnBXSXVXSzZZOFpYblEvYVowMDBMTFVBMVA4N1l3V2FRSWJRTGVPelNhc29GYm5pbkZ3R05FdVdCK0t5MWNMQkRZc1lmZExHQnVYTkRlVmtnUUE3ODJXWWxaNU1lN0RVMWt0Q0U3Qk5jOUlyUVA1YWZDU2g="
}

resource "nifcloud_private_lan" "basic" {
  private_lan_name = "%s"
  cidr_block       = "192.168.100.0/24"
}

resource "nifcloud_route_table" "basic" {
  route {
    cidr_block = "1.1.1.1"
    ip_address = "192.168.100.1"
  }
}

resource "nifcloud_route_table" "upd" {
  route {
    cidr_block = "1.1.1.1"
    ip_address = "192.168.100.1"
  }
}

resource "nifcloud_elb" "basic" {
  elb_name                                    = "%s"
  availability_zone                           = "east-21"
  accounting_type                             = "1"
  network_volume                              = 20
  description                                 = "memo"
  balancing_type                              = 2
  instance_port                               = 3000
  protocol                                    = "HTTPS"
  lb_port                                     = 443
  ssl_certificate_id                          = nifcloud_ssl_certificate.basic.id
  unhealthy_threshold                         = 2
  health_check_target                         = "HTTP:3000"
  health_check_interval                       = 10
  health_check_path                           = "/health"
  health_check_expectation_http_code          = ["200"]
  instances                                   = [nifcloud_instance.basic.instance_id]
  session_stickiness_policy_enable            = true
  session_stickiness_policy_method            = "1"
  session_stickiness_policy_expiration_period = 4
  sorry_page_enable                           = true
  sorry_page_redirect_url                     = "https://example.com"
  route_table_id                              = nifcloud_route_table.basic.id

  network_interface {
    network_name   = nifcloud_private_lan.basic.private_lan_name
    ip_address     = "192.168.100.101"
    is_vip_network = false
  }

  network_interface {
    network_id     = "net-COMMON_GLOBAL"
    is_vip_network = true
  }

  depends_on = [nifcloud_private_lan.basic, nifcloud_route_table.basic, nifcloud_instance.basic, nifcloud_ssl_certificate.basic]
}

resource "nifcloud_ssl_certificate" "basic" {
  certificate = <<EOT
%s
EOT
  key         = <<EOT
%s
EOT
  ca          = <<EOT
%s
EOT
  description = "memo"
}
