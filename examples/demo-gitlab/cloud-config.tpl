#cloud-config
debug:
  verbose: true
cloud_init_modules:
 - migrator
 - seed_random
 - write-files
 - growpart
 - resizefs
 - set_hostname
 - update_hostname
 - update_etc_hosts
 - users-groups
 - ssh
 - runcmd
 - write_files
users:
  - name: "${user_login}"
    sudo: ['ALL=(ALL) NOPASSWD:ALL']
    groups: sudo
    shell: /bin/bash
    ssh_authorized_keys:
    - "${public_key}"
disable_root: true
timezone: "Europe/Moscow"
package_update: false
manage_etc_hosts: localhost
fqdn: "${hostname}"
runcmd:
- apt-get -y update
- apt-get install -y curl openssh-server ca-certificates tzdata perl
- wget https://packages.gitlab.com/install/repositories/gitlab/gitlab-ee/script.deb.sh
- bash ./script.deb.sh
- EXTERNAL_URL=$(curl -s -4 -m 10 https://digitalresistance.dog/myIp) GITLAB_ROOT_PASSWORD="${gitlab_password}" apt-get install gitlab-ee
