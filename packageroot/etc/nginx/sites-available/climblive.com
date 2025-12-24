server {
	listen 443 ssl http2;
	server_name climblive.com;

	client_max_body_size 1M;

	location / {
		root /usr/share/climblive/www;
		try_files $uri /;
		expires 1d;
		add_header Cache-Control "public";
		access_log off;
	}

	location = / {
		root /usr/share/climblive/www;
		try_files /index.html =404;
		add_header Cache-Control "no-store";
		expires 0;

		add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'sha256-jIhoHP5AYEa/rjrf399lCKS/+7hIAc+G1cKDLBSPd7o='; frame-ancestors 'none'; form-action 'none'; base-uri 'self'";
		add_header X-Content-Type-Options "nosniff";
		add_header Referrer-Policy "same-origin";
	}

	ssl_certificate /etc/nginx/ssl/cloudflare/climblive.com/cert.pem;
	ssl_certificate_key /etc/nginx/ssl/cloudflare/climblive.com/privkey.pem;
}

server {
	listen 80;
	server_name climblive.com;
	return 301 https://climblive.com$request_uri;
}
