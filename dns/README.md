docker build -t eh-dns-server .
docker run --rm --name dns-server -v ./etc/bind:/etc/bind -p 53:53/udp -p 53:53/tcp eh-dns-server

arbitary phpmyadmin: 172.20.0.2

zone "example.com" {
    type master;
    file "/etc/bind/db.example.com";
};

zone "0.20.172.in-addr.arpa" {
    type master;
    file "/etc/bind/db.172.20.0";
};


;
; BIND data file for example.com
;
$TTL    604800
@       IN      SOA     ns1.example.com. admin.example.com. (
                        2023091401      ; Serial
                        604800          ; Refresh
                        86400           ; Retry
                        2419200         ; Expire
                        604800 )        ; Negative Cache TTL
;
@       IN      NS      ns1.example.com.
ns1     IN      A       172.20.0.2
www     IN      A       172.20.0.2


~/bind9/etc/bind/db.172.20.0

;
; BIND reverse data file for 172.20.0.0/24
;
$TTL    604800
@       IN      SOA     ns1.example.com. admin.example.com. (
                        2023091401      ; Serial
                        604800          ; Refresh
                        86400           ; Retry
                        2419200         ; Expire
                        604800 )        ; Negative Cache TTL
;
@       IN      NS      ns1.example.com.
10      IN      PTR     example.com.


docker run -d \
  --name bind9 \
  -p 53:53/udp \
  -p 53:53/tcp \
  -v ~/bind9/etc/bind:/etc/bind \
  -v ~/bind9/var/cache/bind:/var/cache/bind \
  --restart=always \
  ubuntu/bind9:latest

docker run --rm \
  --name bind9 \
  -p 53:53/udp \
  -p 53:53/tcp \
  -v ~/bind9/etc/bind:/etc/bind \
  -v ~/bind9/var/cache/bind:/var/cache/bind \
  ubuntu/bind9:latest


docker run --rm \
  --name bind9 \
  -p 53:53/udp \
  -p 53:53/tcp \
  -v ~/bind9/var/cache/bind:/var/cache/bind \
  ubuntu/bind9:latest


  docker run -d \
  --name bind9 \
  -p 53:53/udp \
  -p 53:53/tcp \
  -e BIND9_USER=root \
  -v /epiclabs23/eh/eh-services/dns/etc/bind:/etc/bind \
  -v /epiclabs23/eh/eh-services/dns/var/cache/bind:/var/cache/bind \
  --restart=always \
  ubuntu/bind9:latest
  
  docker run -d \
  --name bind9 \
  -p 53:53/udp \
  -p 53:53/tcp \
  -v /epiclabs23/eh/eh-services/dns/etc/bind:/etc/bind \
  -v /epiclabs23/eh/eh-services/dns/var/cache/bind:/var/cache/bind \
  --restart=always \
  ubuntu/bind9:latest