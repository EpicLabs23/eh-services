### Installation
[Reference 1](https://contabo.com/blog/how-to-setup-your-own-mailserver-with-mailcow/)
[Reference 2](https://docs.mailcow.email/getstarted/prerequisite-system/)

#### Dependecy
Git, Docker

#### Check port availability
```bash
ss -tlpn | grep -E -w '25|80|110|143|443|465|587|993|995|4190'
# or:
netstat -tulpn | grep -E -w '25|80|110|143|443|465|587|993|995|4190'
```
If output is empty, that means, all ports are available.

#### Installation directory
```bash
/root/mailcow-dockerized
```
#### Disable solr
Solr for Mailcow is depricated anyway. So disable it.

```bash
docker stop mailcowdockerized-solr-mailcow-1
```
