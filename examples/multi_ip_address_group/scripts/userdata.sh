#!/bin/bash

configure_multi_ip_address () {
  CIDR=$(python3 -c "import ipaddress, sys; print(ipaddress.IPv4Network('0.0.0.0/${subnet_mask}', strict=False).prefixlen)")
  cat << EOS > /etc/netplan/99-multi-ip-address.yaml
network:
    ethernets:
        ens192:
            dhcp4: no
            dhcp6: no
            addresses: ["${ip_address}/$${CIDR}"]
            routes:
                - to: default
                  via: ${default_gateway}
EOS
  netplan apply
}

configure_multi_ip_address
