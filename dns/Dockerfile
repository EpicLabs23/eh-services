# Use an official BIND image or build from a base image
FROM ubuntu:latest

# Install BIND
RUN apt-get update && apt-get install -y bind9

# Copy configuration files
# COPY named.conf /etc/bind/named.conf
# COPY named.conf.local /etc/bind/named.conf.local
# COPY named.conf.options /etc/bind/named.conf.options
# COPY db.example.com /etc/bind/db.example.com

# Expose the DNS port
EXPOSE 53/udp
EXPOSE 53/tcp

# Start BIND
CMD ["named", "-g"]
