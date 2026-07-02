# assets: build the front-end bundles with Vite (Bootstrap 5, app JS/SCSS).
# The compiled manifest + assets land in public/build and are copied into the
# runtime image below; @vite() reads public/build/manifest.json at request time.
FROM node:22-alpine AS assets

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY vite.config.js ./
COPY resources ./resources
RUN npm run build

# base: PHP-FPM only — used by local docker-compose
FROM php:8.4-fpm-alpine AS base

RUN apk add --no-cache \
    bash \
    curl \
    freetype-dev \
    git \
    icu-dev \
    libjpeg-turbo-dev \
    libpng-dev \
    libzip-dev \
    oniguruma-dev \
    unzip \
    zip

RUN docker-php-ext-configure gd --with-freetype --with-jpeg \
 && docker-php-ext-install \
        bcmath \
        exif \
        gd \
        intl \
        mbstring \
        opcache \
        pdo_mysql \
        zip

COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

WORKDIR /var/www/html

COPY --chown=www-data:www-data . .

RUN git config --global --add safe.directory /var/www/html \
 && composer install --no-interaction --no-scripts --no-dev --optimize-autoloader \
 && chown -R www-data:www-data storage bootstrap/cache

# Vite-built front-end assets (public/build is gitignored, so it is produced
# in the assets stage above rather than copied from the build context).
COPY --from=assets --chown=www-data:www-data /app/public/build ./public/build

COPY docker/php/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["php-fpm"]

# k8s: combined nginx + PHP-FPM — used for Kubernetes deployments
FROM base AS k8s

RUN apk add --no-cache nginx supervisor

COPY docker/nginx/k8s.conf /etc/nginx/http.d/default.conf
COPY docker/supervisord.conf /etc/supervisord.conf

EXPOSE 80

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
