server {
	listen 443 ssl http2;
	server_name __SERVER_NAME__ www.__SERVER_NAME__;

	# Gzip Settings
	include snippets/gzip.conf;

	client_max_body_size 1M;

	set $csp "default-src 'self'; connect-src 'self' clmb.auth.eu-west-1.amazoncognito.com *.fontawesome.com *.sentry.io data:; style-src 'self' https://fonts.googleapis.com 'unsafe-inline'; font-src 'self' https://fonts.gstatic.com; object-src 'none'; frame-ancestors 'none'; form-action 'none'; base-uri 'self'; img-src 'self' data:; report-uri https://o4509937603641344.ingest.de.sentry.io/api/4509937616093264/security/?sentry_key=019099d850441f60cea5d465e217f768";

	location /api {
		rewrite ^/api(.*)$ $1 break;
		proxy_pass http://127.0.0.1:8090;
		proxy_http_version 1.1;
		add_header Cache-Control "no-store";

		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	}

	location /admin {
		root /usr/share/climblive;
		try_files $uri /admin;
		expires 1d;
		add_header Cache-Control "public";
		access_log off;
	}

	location = /admin {
		root /usr/share/climblive/admin;
		try_files /index.html =404;
		add_header Cache-Control "no-store";
		expires 0;
		add_header Content-Security-Policy $csp;
		add_header X-Content-Type-Options "nosniff";
		add_header Referrer-Policy "same-origin";
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
		add_header Cache-Control "no-store";
		expires 0;
		add_header Content-Security-Policy $csp;
		add_header X-Content-Type-Options "nosniff";
		add_header Referrer-Policy "same-origin";
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
		add_header Cache-Control "no-store";
		expires 0;
		add_header Content-Security-Policy $csp;
		add_header X-Content-Type-Options "nosniff";
		add_header Referrer-Policy "same-origin";
	}

	include /etc/nginx/options-ssl.conf;

	ssl_certificate /etc/nginx/ssl/cloudflare/climblive.app/cert.pem;
	ssl_certificate_key /etc/nginx/ssl/cloudflare/climblive.app/privkey.pem;
}

server {
	listen 80;
	server_name __SERVER_NAME__ www.__SERVER_NAME__;
	return 301 https://__SERVER_NAME__$request_uri;
}
