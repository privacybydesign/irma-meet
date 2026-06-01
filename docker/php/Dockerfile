# assets: Node 20 + Vite — produces public/build/ for inclusion in `base`.
FROM node:20-alpine AS assets

WORKDIR /app

COPY package.json package-lock.json* ./
RUN npm ci

COPY resources ./resources
COPY vite.config.js ./

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
    libpq \
    libzip-dev \
    oniguruma-dev \
    unzip \
    zip

# postgresql-dev ships the headers needed to compile pdo_pgsql, but the
# resulting extension only needs libpq at runtime. Stage it via a virtual
# package and delete it once compilation finishes to keep the image small.
RUN apk add --no-cache --virtual .build-deps postgresql-dev \
 && docker-php-ext-configure gd --with-freetype --with-jpeg \
 && docker-php-ext-install \
        bcmath \
        exif \
        gd \
        intl \
        mbstring \
        opcache \
        pdo_pgsql \
        zip \
 && apk del .build-deps

COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

WORKDIR /var/www/html

COPY --chown=www-data:www-data . .

# Vite output from the assets stage; not committed.
COPY --from=assets --chown=www-data:www-data /app/public/build ./public/build

RUN git config --global --add safe.directory /var/www/html \
 && composer install --no-interaction --no-scripts --no-dev --optimize-autoloader \
 && chown -R www-data:www-data storage bootstrap/cache

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
