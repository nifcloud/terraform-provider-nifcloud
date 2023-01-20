#!/bin/bash

configure_ssh_port () {
  sed -i 's/^#*Port [0-9]*/Port ${custom_ssh_port}/' /etc/ssh/sshd_config
}

configure_ssh_port
