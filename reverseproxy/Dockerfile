FROM nginx:alpine

RUN mkdir -p /etc/letsencrypt/live/halspals.co.uk/
COPY nginx.conf /etc/nginx/nginx.conf
COPY /etc/letsencrypt/live/halspals.co.uk/fullchain.pem /etc/letsencrypt/live/halspals.co.uk/
COPY /etc/letsencrypt/live/halspals.co.uk/privkey.pem /etc/letsencrypt/live/halspals.co.uk/
COPY /etc/letsencrypt/options-ssl-nginx.conf /etc/letsencrypt/
COPY /etc/letsencrypt/ssl-dhparams.pem /etc/letsencrypt/
