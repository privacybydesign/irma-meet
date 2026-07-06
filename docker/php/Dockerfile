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
    sqlite-dev \
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
        pdo_sqlite \
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

# test: run the PHPUnit suite as a build-time gate. docker-image.yml builds this
# image on every push/PR with target k8s, and k8s copies the stamp produced here,
# so BuildKit is forced to build this stage and a failing test aborts the whole
# image build -> the CI job goes red. This is how PHPUnit actually runs in CI
# today (issue #27) without a dedicated GitHub Actions test workflow.
FROM base AS test

# The base stage installed with --no-dev, so phpunit/mockery/faker/collision are
# absent; restore the dev dependencies just for this throwaway stage.
RUN composer install --no-interaction --no-scripts --optimize-autoloader

# View-rendering feature tests (e.g. GET /) resolve @vite(...), which reads the
# compiled manifest — pull it in from the assets stage.
COPY --from=assets --chown=www-data:www-data /app/public/build ./public/build

# A minimal, self-contained .env is enough: phpunit.xml supplies the DB (sqlite
# :memory:), cache, mail and session drivers; the app just needs a valid APP_KEY.
# (Kept independent of .env.example, whose placeholder values are unquoted.)
RUN printf 'APP_NAME=irma-meet\nAPP_ENV=testing\nAPP_KEY=\nAPP_DEBUG=true\nAPP_URL=http://localhost\n' > .env \
 && php artisan key:generate \
 && php artisan test \
 && touch /tmp/tests-passed

# k8s: combined nginx + PHP-FPM — used for Kubernetes deployments
FROM base AS k8s

# Depend on the test stage's stamp so BuildKit builds (and thus runs) the test
# stage before producing the deployable image. The stamp itself is inert.
COPY --from=test /tmp/tests-passed /tmp/tests-passed

RUN apk add --no-cache nginx supervisor

COPY docker/nginx/k8s.conf /etc/nginx/http.d/default.conf
COPY docker/supervisord.conf /etc/supervisord.conf

EXPOSE 80

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
