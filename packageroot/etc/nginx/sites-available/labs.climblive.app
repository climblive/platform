server {
	listen 443 ssl;
	server_name labs.climblive.app;

	http2 on;

	# Gzip Settings
	include snippets/gzip.conf;

	client_max_body_size 1M;

	location /api {
		rewrite ^/api(.*)$ $1 break;
		proxy_pass http://127.0.0.1:8090;
		proxy_http_version 1.1;
		add_header Cache-Control "no-store";

		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	}

	location /scoreboard {
		root /usr/share/climblive;
		try_files $uri /scoreboard;
		expires 1d;
		add_header Cache-Control "public";
		access_log off;
	}

	location = /scoreboard {
		root /usr/share/climblive/scoreboard;
		try_files /index.html =404;
	}

	location / {
		root /usr/share/climblive/scorecard;
		try_files $uri /;
		expires 1d;
		add_header Cache-Control "public";
		access_log off;
	}

	location = / {
		root /usr/share/climblive/scorecard;
		try_files /index.html =404;
	}

	add_header Content-Security-Policy "default-src 'self'; connect-src 'self' api.labs.climblive.app data:; style-src 'self' https://fonts.googleapis.com 'unsafe-inline'; font-src 'self' https://fonts.gstatic.com; object-src 'none'; frame-ancestors 'none'; form-action 'none'; base-uri 'self'";
	add_header X-Content-Type-Options "nosniff";
	add_header Referrer-Policy "same-origin";

	include /etc/nginx/options-ssl.conf;

	ssl_certificate /etc/letsencrypt/live/labs.climblive.app/fullchain.pem;
	ssl_certificate_key /etc/letsencrypt/live/labs.climblive.app/privkey.pem;
}