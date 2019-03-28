#!/bin/sh
echo "
worker_processes 1;

events { worker_connections 1024; }

http {
    server {
            server_name hydra.$1;
            location / {
                proxy_pass http://localhost:4445;
        }


    listen 443 ssl; # managed by Certbot
        ssl_certificate /etc/letsencrypt/live/client.api.$1/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/client.api.$1/privkey.pem; # managed by Certbot
        include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot

    }

        server {
        server_name oauth.$1;
        location / {
            proxy_pass http://localhost:4444;
            }


        listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/client.api.$1/fullchain.pem; # managed by Certbot
        ssl_certificate_key /etc/letsencrypt/live/client.api.$1/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
        ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot


}

    server {
            server_name account.$1;
            location / {
                proxy_pass http://localhost:4001;
        }


    listen 443 ssl; # managed by Certbot
        ssl_certificate /etc/letsencrypt/live/client.api.$1/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/client.api.$1/privkey.pem; # managed by Certbot
        include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot



    }

        server {
        server_name portainer.$1;
        location / {
            proxy_pass http://localhost:9000;
            }
        listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/client.api.$1/fullchain.pem; # managed by Certbot
        ssl_certificate_key /etc/letsencrypt/live/client.api.$1/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
        ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot




}
    server {
            server_name adminer.$1;
            location / {
                proxy_pass http://localhost:8080;
        }


    listen 443 ssl; # managed by Certbot
        ssl_certificate /etc/letsencrypt/live/client.api.$1/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/client.api.$1/privkey.pem; # managed by Certbot
        include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot




    }
        server {
        server_name robot.$1;
        location / {
            proxy_pass http://localhost:8083;
                proxy_set_header Upgrade \$http_upgrade;
            proxy_set_header Connection \"upgrade\";
                proxy_set_header Origin \"\";
        }
}
    server {
            server_name client.api.$1;
            location / {
                proxy_pass http://localhost:9100;
        }


    listen 443 ssl; # managed by Certbot
        ssl_certificate /etc/letsencrypt/live/client.api.$1/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/client.api.$1/privkey.pem; # managed by Certbot
        include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot




    }
        server {
        server_name portainer.test.$1;
        location / {
            proxy_pass http://localhost:11000;
            }
        listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/client.api.$1/fullchain.pem; # managed by Certbot
        ssl_certificate_key /etc/letsencrypt/live/client.api.$1/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
        ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot




}
    server {
            server_name adminer.test.$1;
            location / {
                proxy_pass http://localhost:10000;
        }


    listen 443 ssl; # managed by Certbot
        ssl_certificate /etc/letsencrypt/live/client.api.$1/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/client.api.$1/privkey.pem; # managed by Certbot
        include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot




    }
        server {
        server_name robot.test.$1;
        location / {
            proxy_pass http://localhost:12000;
            }



    }
        server {
        server_name client.api.test.$1;
        location / {
            proxy_pass http://localhost:13000;
            }


        listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/client.api.$1/fullchain.pem; # managed by Certbot
        ssl_certificate_key /etc/letsencrypt/live/client.api.$1/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
        ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot




}
server {
        server_name $1;
        location / {
                proxy_pass http://localhost:5030;
        }

        listen 443 ssl;

    ssl_certificate /etc/letsencrypt/live/$1/fullchain.pem; # managed by Certbot
        ssl_certificate_key /etc/letsencrypt/live/$1/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
        ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot

}

    server {
        if (\$host = adminer.$1) {
        return 301 https://\$host\$request_uri;
    } # managed by Certbot


            server_name adminer.$1;
        listen 80;
    return 404; # managed by Certbot


    }
        server {
    if (\$host = client.api.$1) {
            return 301 https://\$host\$request_uri;
        } # managed by Certbot


        server_name client.api.$1;
    listen 80;
        return 404; # managed by Certbot


}
    server {
        if (\$host = portainer.$1) {
        return 301 https://\$host\$request_uri;
    } # managed by Certbot


            server_name portainer.$1;
        listen 80;
    return 404; # managed by Certbot


    }

        server {
    if (\$host = adminer.test.$1) {
            return 301 https://\$host\$request_uri;
        } # managed by Certbot
        server_name adminer.test.$1;
    listen 80;
        return 404; # managed by Certbot


}
    server {
        if (\$host = client.api.test.$1) {
        return 301 https://\$host\$request_uri;
    } # managed by Certbot
    
    
            server_name client.api.test.$1;
        listen 80;
    return 404; # managed by Certbot
    
    
    }
        server {
    if (\$host = portainer.test.$1) {
            return 301 https://\$host\$request_uri;
        } # managed by Certbot


        server_name portainer.test.$1;
    listen 80;
        return 404; # managed by Certbot


}
    server {
        if (\$host = robot.test.$1) {
        return 301 https://\$host\$request_uri;
    } # managed by Certbot
    
    
            server_name robot.test.$1;
        listen 80;
    return 404; # managed by Certbot
    
    
    }
        server {
    if (\$host = account.$1) {
            return 301 https://\$host\$request_uri;
        } # managed by Certbot


        server_name account.$1;
    listen 80;
        return 404; # managed by Certbot


}
    server {
        if (\$host = oauth.$1) {
        return 301 https://\$host\$request_uri;
    } # managed by Certbot
    
            server_name oauth.$1;
        listen 80;
    return 404; # managed by Certbot


    }
        server {
    if (\$host = hydra.$1) {
            return 301 https://\$host\$request_uri;
        } # managed by Certbot


        server_name hydra.$1;
    listen 80;
        return 404; # managed by Certbot


}}
    " > /etc/nginx/nginx.conf
