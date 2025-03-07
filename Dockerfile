FROM postgres:17

COPY init.sql /docker-entrypoint-initdb.d/init.sql

RUN sed -i 's/psql "$@"/psql "$@" -v ON_ERROR_STOP=0/g' /usr/local/bin/docker-entrypoint.sh