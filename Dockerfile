FROM alpine

RUN apk update && \
    apk add tzdata wget && \
    cp /usr/share/zoneinfo/Asia/Singapore /etc/localtime && \
    echo "Asia/Singapore" >  /etc/timezone && \
    addgroup -g 1000 goschedule && adduser -D -G goschedule -u 1000 goschedule && \
    mkdir /app && cd /app && chown goschedule:goschedule /app && \
    wget -O gosu https://github.com/tianon/gosu/releases/download/1.14/gosu-amd64 && chmod +x gosu && mv gosu /usr/bin/ && \
    wget https://github.com/jasonjoo2010/goschedule-console/releases/download/v1.1.0/goschedule-console-v1.1.0-linux-x86_64.tar.gz && \
    tar zxf goschedule-console-*.tar.gz && rm goschedule-console-*.tar.gz && \
    mv goschedule-console* /usr/bin/goschedule-console && \
    apk del tzdata wget

WORKDIR /app
COPY docker-entrypoint.sh /usr/bin/
EXPOSE 8000
ENTRYPOINT ["docker-entrypoint.sh"]
