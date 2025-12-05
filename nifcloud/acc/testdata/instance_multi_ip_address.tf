provider "nifcloud" {
  region = "jp-east-2"
}

resource "nifcloud_instance" "basic" {
  instance_id             = "%s"
  description             = "memo"
  availability_zone       = "east-21"
  accounting_type         = "2"
  image_id                = "354"
  instance_type           = "c2-small"
  key_name                = nifcloud_key_pair.basic.key_name
  security_group          = nifcloud_security_group.basic.group_name
  user_data               = "#!/bin/bash"
  
  multi_ip_address_configuration_user_data = <<-EOT
    #!/bin/bash
    
    configure_private_ip_address () {
      CIDR=$(python3 -c "import ipaddress, sys; print(ipaddress.IPv4Network('0.0.0.0/${nifcloud_multi_ip_address_group.basic.subnet_mask}', strict=False).prefixlen)")
      cat << EOS > /etc/netplan/99-multi-ip-address.yaml
    network:
        ethernets:
            ens192:
                dhcp4: no
                dhcp6: no
                addresses: ["${nifcloud_multi_ip_address_group.basic.ip_addresses[0]}/$${CIDR}"]
                routes:
                    - to: default
                      via: ${nifcloud_multi_ip_address_group.basic.default_gateway}
    EOS
      netplan apply
    }
        
    configure_private_ip_address
  EOT

  network_interface {
    network_id = "net-COMMON_PRIVATE"
  }

  network_interface {
    network_id = "net-MULTI_IP_ADDRESS"
    multi_ip_address_group_id = nifcloud_multi_ip_address_group.basic.id
  }

  depends_on = [nifcloud_key_pair.basic, nifcloud_security_group.basic, nifcloud_multi_ip_address_group.basic]
}

resource "nifcloud_security_group" "basic" {
  group_name        = "%s"
  availability_zone = "east-21"
}

resource "nifcloud_key_pair" "basic" {
  key_name   = "%s"
  public_key = "c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFEeFVVSmtIWFFvalVmeGphT3dQNVJmMjhOTVRFSjJFblBQdFk0b1NkZFBpRllnMWVDTGFNU08wV25nZVIrVk5sU215am1qU2xRWjBsc1BkcHZjWnY0KzZiMDlLUUZlT3NxakdjNE9Ga1o2MTZyTEI3UmdzblZnSXl3QmtIZ2lsMVQzbFRwRHVtYVk2TFFaRjRiaVpTNkNyaFdYeVhiSjFUVmYyZ0hIYXZPdi9WSS9ITjhIejlnSDg5Q0xWRVFOWFVQbXdjbC83ZE4yMXE4QnhNVkpGNW1sSW1RcGxwTjFKQVRwdnBXSXVXSzZZOFpYblEvYVowMDBMTFVBMVA4N1l3V2FRSWJRTGVPelNhc29GYm5pbkZ3R05FdVdCK0t5MWNMQkRZc1lmZExHQnVYTkRlVmtnUUE3ODJXWWxaNU1lN0RVMWt0Q0U3Qk5jOUlyUVA1YWZDU2g="
}

resource "nifcloud_multi_ip_address_group" "basic" {
  name              = "%s"
  availability_zone = "east-21"
  ip_address_count  = 1
}
