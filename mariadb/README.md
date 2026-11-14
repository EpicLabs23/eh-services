## Installation
See [https://github.com/EpicLabs23/eh-services](https://github.com/EpicLabs23/eh-services)
## Usage
phpMyAdmin is no longer published on a host port. It is served over TLS on `https://<account primary domain>:2329` by the host nginx (EHM generates the vhost). For direct local debugging, use the container address: http://172.1.0.5/