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
users:
  - name: debian
    sudo: ['ALL=(ALL) NOPASSWD:ALL']
    groups: sudo
    shell: /bin/bash
chpasswd:
  list:
    - "debian:${vm_password}"
  expire: False
fqdn: "${hostname}"
runcmd:
- apt-get -y update
- apt-get -y install nginx
- echo '<h1>Hello, World!</h1><code>Node ${hostname}</code>' > /var/www/html/index.nginx-debian.html
