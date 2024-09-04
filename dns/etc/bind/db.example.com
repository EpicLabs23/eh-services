$TTL 86400
@   IN  SOA ns1.example.com. admin.example.com. (
            2024090401 ; Serial
            3600       ; Refresh
            1800       ; Retry
            1209600    ; Expire
            86400 )    ; Minimum TTL

@       IN  NS  ns1.example.com.
ns1     IN  A   192.168.1.100
www     IN  A   192.168.1.101
