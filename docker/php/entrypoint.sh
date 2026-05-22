#!/bin/sh
set -e

cd /var/www/html

# Ensure writable directories exist
mkdir -p storage/framework/{cache,sessions,views} storage/logs bootstrap/cache
chown -R www-data:www-data storage bootstrap/cache

# Run deferred Composer scripts now that .env is present
composer run-script post-autoload-dump --no-interaction 2>/dev/null || true

exec "$@"
