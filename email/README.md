### Referrence Links
[Email Components](https://www.digitalocean.com/community/tutorials/why-you-may-not-want-to-run-your-own-mail-server)
[Install and configure postfix on ubuntu 22](https://www.digitalocean.com/community/tutorials/how-to-install-and-configure-postfix-on-ubuntu-22-04)
[DNS Black list check](https://www.dnsbl.info/dnsbl-database-check.php)
[MX Toolbox](https://mxtoolbox.com/SuperTool.aspx)
[Email Tester](https://www.mail-tester.com/)
[Install and Configure Mail Server on Ubuntu 22.04 with Postfix, Dovecot and Roundcube](https://www.tuxnoob.com/posts/Install-and-Configure-Mail-Server-ubuntu-part1/)

### Cheatsheet
Send a mail:
```bash
cat ~/test_message | s-nail -s 'Test email subject line' -r admin@ehm23.com test-szf7h907a@srv1.mail-tester.com
```
s-nail commands:
```bash
# Check sent items
? file +sent
```

Postfix commands:
```bash
# Check postfix logs
sudo tail -f /var/log/mail.log
#Mail Queue Status
sudo postfix queue
sudo postqueue -p
# Flush queue
sudo postfix flush
```

### Email Components
- MTA ( Mail Transfer Agent ), resposible for sending email to another MTA i.e. another mail server. 