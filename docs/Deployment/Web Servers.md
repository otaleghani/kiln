---
title: Deploy on Web Servers
description: Configuration guides for hosting your Kiln site on your own VPS using Nginx, Apache, or Caddy.
---
# Deploy on Web Servers

If you have your own Virtual Private Server (VPS) or a Raspberry Pi, you can host your Kiln site using any standard web server.

The process is simple:
1.  Run `./kiln generate` on your local machine (or on the server).
2.  Upload the contents of the `public` folder to your server's web root.
3.  Configure your web server to serve static files.

Below are the configuration snippets for the most popular web servers.

## Nginx

Add this `server` block to your Nginx configuration (usually in `/etc/nginx/sites-available/default`).

```nginx
server {
    listen 80;
    server_name example.com;
    
    # Point this to where you uploaded your 'public' folder
    root /var/www/kiln/public;
    index index.html;

    location / {
        # Try to serve the file directly, then as a directory (index.html), then 404
        try_files $uri $uri/ =404;
    }

    # Optional: Cache static assets for better performance
    location ~* \.(css|js|png|jpg|jpeg|gif|ico)$ {
        expires 30d;
        add_header Cache-Control "public, no-transform";
    }
}
```

## Apache

Ensure your VirtualHost points to the public directory. You may also need to enable the `mod_rewrite` module if you plan to do complex routing, but for standard Kiln sites, a basic config works.

```apache
<VirtualHost *:80>
    ServerName example.com
    DocumentRoot /var/www/kiln/public

    <Directory /var/www/kiln/public>
        Options Indexes FollowSymLinks
        AllowOverride All
        Require all granted
    </Directory>

    ErrorLog ${APACHE_LOG_DIR}/error.log
    CustomLog ${APACHE_LOG_DIR}/access.log combined
</VirtualHost>
```

## Caddy

Caddy is arguably the easiest server to configure for static sites as it handles HTTPS automatically. Create a `Caddyfile` in your directory:
```Caddyfile
example.com {
    # Point this to where you uploaded your 'public' folder
    root * /var/www/kiln/public
    
    # Try the exact path first, then try appending .html 
	try_files {path} {path}.html
    
    # Enable static file server
    file_server
    
    # Compress responses for speed
    encode gzip
}
```

If you want to deploy under a directory, you can use the following `Caddyfile` as an example:
```Caddyfile
example.com {
    handle_path /subpath/* {
	    # Point this to where you uploaded your 'public' folder
		root * /var/www/fictitiousentrycom/subpath
		
		# Try the exact path first, then try appending .html
		try_files {path} {path}.html
		
		file_server
	}
}
```