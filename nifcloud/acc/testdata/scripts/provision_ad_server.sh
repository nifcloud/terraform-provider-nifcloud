#!/bin/bash

set -euo pipefail

function change_apt_repository () {
  sed -i.bak -e "s%http://[^ ]\+%http://ftp.riken.go.jp/Linux/ubuntu/%g" /etc/apt/sources.list
}

function disable_auto_update () {
  echo 'APT::Periodic::Enable "0";' | tee /etc/apt/apt.conf.d/99disable-periodic
}

function configure_ad_server () {
  HOST_NAME="AD01"
  DOMAIN="TFACC"
  REALM="TFACC.LOCAL"
  HOST_IP="$(ip a show ens224 | grep -o 'inet [0-9]\+\.[0-9]\+\.[0-9]\+\.[0-9]\+' | grep -o [0-9].*)" 
  echo $${HOST_IP} > /var/log/ip

  apt-get update
  DEBIAN_FRONTEND=noninteractive apt-get -y install slapd ldap-utils libnss-ldap samba smbldap-tools smbclient krb5-user winbind

  rm /etc/samba/smb.conf
  samba-tool domain provision --use-rfc2307 --realm=$${REALM} --server-role=dc --dns-backend=SAMBA_INTERNAL --domain=$${DOMAIN} --host-name=$${HOST_NAME} --host-ip=$${HOST_IP} --adminpass=tfaccpass+555 --option="interfaces=lo ens224"

  cp -p /var/lib/samba/private/krb5.conf /etc/krb5.conf
  cat << EOF >> /etc/krb5.conf
[realms]
    $${REALM} = {
        kdc = $${HOST_NAME,,}.$${REALM,,}
        admin_server = $${HOST_NAME,,}.$${REALM,,}
    }
[domain_realm]
    .$${REALM,,} = $${REALM}
    $${REALM,,} = $${REALM}
EOF

  rm /etc/resolv.conf
  cat << EOF > /etc/resolv.conf
nameserver 127.0.0.1
nameserver 8.8.8.8
domain $${REALM,,}
EOF

  systemctl stop slapd
  systemctl disable slapd
  systemctl stop systemd-resolved
  systemctl disable systemd-resolved
  systemctl stop unbound
  systemctl disable unbound
  systemctl stop winbind
  systemctl disable winbind
  systemctl stop nmbd.service
  systemctl disable nmbd.service
  systemctl stop smbd.service
  systemctl disable smbd.service
  systemctl unmask samba-ad-dc.service
  systemctl enable samba-ad-dc.service
  systemctl restart samba-ad-dc.service
}

disable_auto_update
change_apt_repository
configure_ad_server
