FROM adminer:latest

USER root

RUN apk add autoconf gcc g++ make libffi-dev openssl-dev \
        && pecl install mongodb \
        && echo "extension=mongodb.so" > /usr/local/etc/php/conf.d/docker-php-ext-mongodb.ini

RUN sed -i "s|{\$h->connect(\"mongodb://\$N\",\$yf);return|{|" adminer.php \
        && sed -i "s|lang(22);}|}|" adminer.php

USER adminer
