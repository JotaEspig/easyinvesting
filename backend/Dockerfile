FROM golang:1.24.3 as builder

WORKDIR /app
COPY . .

RUN go build -o app .

FROM debian:bullseye

RUN apt-get update && \
    apt-get install -y apache2 openssl && \
    a2enmod ssl proxy proxy_http headers && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/app /usr/local/bin/app

COPY apache/backend.conf /etc/apache2/sites-available/go-backend.conf
COPY apache/ssl /etc/apache2/ssl

RUN a2ensite go-backend.conf && a2dissite 000-default && \
    mkdir -p /run/apache2

EXPOSE 443

ENV EASYINVESTING_PORT=8080
ENV EASYINVESTING_DB_STRING=/app/easyinvesting.db
ENV ENV=PRODUCTION

# Change these to your actual secrets
ENV EASYINVESTING_SECRET_KEY=default_secret_key
ENV BRAPI_TOKEN=default_brapi_token

CMD ["/bin/bash", "-c", "/usr/local/bin/app & apache2ctl -D FOREGROUND"]
