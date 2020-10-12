FROM aa

VOLUME /var/lib/redis
WORKDIR /var/lib/redis
COPY ./docker-healthcheck /usr/bin/
COPY ./redis.conf /etc/redis.conf
RUN chmod +x /usr/bin/docker-healthcheck \
    && chown redis:redis /etc/redis.conf

HEALTHCHECK CMD ["docker-healthcheck"]
USER redis
EXPOSE 6379
CMD ["redis-server", "/etc/redis.conf"]
