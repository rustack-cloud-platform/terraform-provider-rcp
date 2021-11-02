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
write_files:
- path: /etc/nginx/sites-enabled/default
  permissions: "0644"
  content: |
   upstream demo {
     ${balancer_upstream}
   }
   server {
     listen 80 default_server;
     root /var/www/html;
     server_name _;
     location / {
       proxy_pass http://demo;
       proxy_read_timeout     120;
       proxy_connect_timeout  120;
     }
   }
