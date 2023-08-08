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

resource "nifcloud_elb" "basic" {
  elb_name          = "%s"
  availability_zone = "east-21"
  instance_port     = 3000
  protocol          = "HTTP"
  lb_port           = 80

  network_interface {
    network_name   = nifcloud_private_lan.basic.private_lan_name
    ip_address     = "192.168.100.101"
    is_vip_network = false
    system_ip_addresses {
      system_ip_address = "192.168.100.102"
    }
    system_ip_addresses {
      system_ip_address = "192.168.100.103"
    }
  }

  network_interface {
    network_id     = "net-COMMON_GLOBAL"
    is_vip_network = true
  }

  depends_on = [nifcloud_private_lan.basic]
}

resource "nifcloud_elb_listener" "basic" {
  elb_id                                      = nifcloud_elb.basic.id
  description                                 = "memo-upd"
  balancing_type                              = 1
  instance_port                               = 3001
  protocol                                    = "HTTP"
  lb_port                                     = 8080
  unhealthy_threshold                         = 3
  health_check_target                         = "HTTP:3001"
  health_check_interval                       = 11
  health_check_path                           = "/health-upd"
  health_check_expectation_http_code          = ["3xx"]
  instances                                   = [nifcloud_instance.upd.instance_id]
  session_stickiness_policy_enable            = true
  session_stickiness_policy_method            = "2"
  session_stickiness_policy_expiration_period = 5
  sorry_page_enable                           = true
  sorry_page_redirect_url                     = "http://example.com"

  depends_on = [nifcloud_elb.basic, nifcloud_instance.upd]
}

resource "tls_private_key" "basic" {
  algorithm = "RSA"
}

resource "tls_self_signed_cert" "basic" {
  private_key_pem       = tls_private_key.basic.private_key_pem
  validity_period_hours = 3
  dns_names             = ["example.com"]
  allowed_uses          = ["client_auth"]

  subject {
    common_name  = "example.com"
    organization = "ACME Examples, Inc"
  }
}

resource "nifcloud_ssl_certificate" "basic" {
  certificate = tls_self_signed_cert.basic.cert_pem
  key         = tls_private_key.basic.private_key_pem
}
