FROM ubuntu/bind9:latest

ENV DEBIAN_FRONTEND=noninteractive

# Install dependencies
RUN apt-get update && apt-get install -y curl supervisor \
    && rm -rf /var/lib/apt/lists/*

# Install bind9-api
RUN curl -sSL https://raw.githubusercontent.com/EpicLabs23/bind9-api/refs/heads/main/install.sh | bash -s 0.1.0 no-systemd

# Create supervisord config
RUN mkdir -p /etc/supervisor/conf.d
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Expose DNS + API ports
EXPOSE 53/udp 53/tcp 8053/tcp

# Start supervisor (which runs both bind9 & bind9-api)
CMD ["/usr/bin/supervisord", "-n", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
