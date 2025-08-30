### Pull the BIND9 docker image before disabling built-in DNS server

```bash
docker pull ubuntu/bind9:latest
```

Git clone or pull: https://github.com/EpicLabs23/eh-services.git /epiclabs23/eh/eh-services/dns

### Disable Ubuntu built-in DNS resolver

1. Stop the service `service systemd-resolved stop`
2. Backup existing `/etc/resolv.conf` just for safety. `mv /etc/resolv.conf /etc/resolv.conf.backup`
3. Disable the service from startup `systemctl disable systemd-resolved`
4. Confirm if port 53 is free: `sudo lsof -i :53` . this should return nothing.

### Installation:

1. `cd /epiclabs23/eh/eh-services/dns`
2. Start docker container:

```bash
docker compose build
docker compose up -d

# docker run -d \
# --name bind9 \
# -p 53:53/udp \
# -p 53:53/tcp \
# -e TZ=Asia/Dhaka \
# -e BIND9_USER=root \
# -v /epiclabs23/eh/eh-services/dns/etc/bind:/etc/bind \
# -v /epiclabs23/eh/eh-services/dns/var/cache/bind:/var/cache/bind \
# -v /epiclabs23/eh/eh-services/dns/var/lib/bind:/var/lib/bind \
# --restart=always \
# ubuntu/bind9:latest
```

3. Configure Your Local Machine to Use the DNS Server: `vim /etc/resolv.conf` then add following content

```bash
nameserver 127.0.0.1
nameserver 8.8.8.8
nameserver 45.125.222.158
search .
```

4. Test DNS server:

```bash
dig @127.0.0.1 www.example.local
```

or

```bash
nslookup www.example.local 127.0.0.1

```

5. Managing DNS server:

```bash
docker exec bind9 named-checkzone ehm23.com /etc/bind/zones/db.ehm23.com
docker exec bind9 named-checkzone 57.169.49.103.in-addr.arpa /etc/bind/zones/db.57.169.49.103
docker exec bind9 rndc reload
docker stop bind9
docker start bind9
docker logs bind9
```

6. Settings, forward DNS etc are available in: `./etc/bind/named.conf.options`

### Adding a new zone file:

1. Copy the `db.example.local` file as template, make changes accordingly.
2. For any subdomain add an 'A' entry on this very same file.
3. Update `named.conf.local` with the new zone file.

### Debug

Start built-in DNS: `service systemd-resolved start`

Check status of built-in DNS: `service systemd-resolved status`

Enable the built-in DSN in startup `systemctl enable systemd-resolved`

Copy back the backed up config file: `cp /etc/resolv.conf.backup /etc/resolv.conf`
