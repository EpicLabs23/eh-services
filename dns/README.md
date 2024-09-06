### Installation:
1. `cd /epiclabs23/eh/eh-services/dns`
2. Start docker container:
```bash
docker run -d \
--name bind9 \
-p 53:53/udp \
-p 53:53/tcp \
-e BIND9_USER=root \
-v /epiclabs23/eh/eh-services/dns/etc/bind:/etc/bind \
-v /epiclabs23/eh/eh-services/dns/var/cache/bind:/var/cache/bind \
--restart=always \
ubuntu/bind9:latest
```
3. Configure Your Local Machine to Use the DNS Server:  `vim /etc/resolv.conf` then add `nameserver 127.0.0.1
` at the top.
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
docker stop bind9
docker start bind9
docker logs bind9
```
6. Settings, forward DNS etc are available in: `./etc/bind/named.conf.options`

### Adding a new zone file:
1. Copy the `db.example.local` file as template, make changes accordingly. for any subdomain add an 'A' entry on this very same file.
2. update `named.conf.local` with the new zone file.