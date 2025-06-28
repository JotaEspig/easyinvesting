FROM golang:latest AS builder

WORKDIR /app
COPY . .

RUN go build -o app .

FROM debian:bookworm

RUN apt-get update && \
    apt-get install -y apache2 openssl && \
    a2enmod ssl proxy proxy_http headers && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/app /usr/local/bin/app

COPY apache/backend.conf /etc/apache2/sites-available/go-backend.conf
COPY apache/ssl /etc/apache2/ssl

RUN echo "ServerName localhost" >> /etc/apache2/apache2.conf && \
    a2ensite go-backend.conf && a2dissite 000-default && \
    mkdir -p /run/apache2

EXPOSE 443

RUN mkdir -p /data && chown www-data:www-data /data
ENV EASYINVESTING_DB_STRING=/data/easyinvesting.db
ENV EASYINVESTING_PORT=8080
ENV ENV=PRODUCTION

CMD ["/bin/bash", "-c", "/usr/local/bin/app & apache2ctl -D FOREGROUND"]
