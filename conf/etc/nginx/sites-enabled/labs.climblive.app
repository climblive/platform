server {
	listen 443 ssl;
	server_name labs.climblive.app;

	http2 on;

	# Gzip Settings
	include snippets/gzip.conf;

	client_max_body_size 1M;

	location /scoreboard {
		root /opt/climblive/labs;
		try_files $uri /scoreboard;
		expires 1d;
		add_header Cache-Control "public";
		access_log off;
	}

	location = /scoreboard {
		root /opt/climblive/labs/scoreboard;
		try_files /index.html =404;
	}

	location / {
		root /opt/climblive/labs/scorecard;
		try_files $uri /;
		expires 1d;
		add_header Cache-Control "public";
		access_log off;
	}

	location = / {
		root /opt/climblive/labs/scorecard;
		try_files /index.html =404;
	}

	include /etc/nginx/options-ssl.conf;

	ssl_certificate /etc/letsencrypt/live/labs.climblive.app/fullchain.pem;
	ssl_certificate_key /etc/letsencrypt/live/labs.climblive.app/privkey.pem;
}

server {
	listen 443 ssl;
	server_name api.labs.climblive.app;

	http2 on;

	# Gzip Settings
	include snippets/gzip.conf;

	client_max_body_size 1M;

	location / {
		proxy_pass http://127.0.0.1:8090;
		proxy_http_version 1.1;
		add_header Cache-Control "no-store";
	}

	include /etc/nginx/options-ssl.conf;

	ssl_certificate /etc/letsencrypt/live/labs.climblive.app/fullchain.pem;
	ssl_certificate_key /etc/letsencrypt/live/labs.climblive.app/privkey.pem;
}
