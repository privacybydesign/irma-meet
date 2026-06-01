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
    postgresql-dev \
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
        pdo_pgsql \
        zip

COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

WORKDIR /var/www/html

COPY --chown=www-data:www-data . .

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
