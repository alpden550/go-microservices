FROM caddy:alpine

COPY Caddyfile /etc/caddy/Caddyfile

COPY errors/404.html .
