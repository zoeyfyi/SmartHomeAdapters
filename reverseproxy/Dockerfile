FROM nginx:alpine

RUN mkdir -p /etc/letsencrypt/live/halspals.co.uk/
COPY nginx.conf /etc/nginx/nginx.conf
COPY fullchain.pem /etc/letsencrypt/live/halspals.co.uk/
COPY privkey.pem /etc/letsencrypt/live/halspals.co.uk/
COPY options-ssl-nginx.conf /etc/letsencrypt/
COPY ssl-dhparams.pem /etc/letsencrypt/
