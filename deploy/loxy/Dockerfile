FROM debian:latest

ENV APP_HOME /home/loxy

RUN apt-get -qq update \
    && apt-get install -y \
       wget \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd -r loxy --gid=1000 \
      && useradd -r -m -g loxy --uid=1000 loxy

WORKDIR ${APP_HOME}

## <[ Kaigara
ENV KAIGARA_VERSION v0.0.2
RUN wget --quiet https://github.com/mod/kaigara/releases/download/$KAIGARA_VERSION/kaigara-linux-amd64-$KAIGARA_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf kaigara-linux-amd64-$KAIGARA_VERSION.tar.gz

COPY operations /opt/provision/operations
COPY resources /opt/provision/resources
COPY metadata.yml /opt/provision
## Kaigara ]>

VOLUME ["${APP_HOME}"]

COPY loxybeat /usr/local/bin/
COPY run-loxybeat.sh /usr/local/bin/

ENTRYPOINT ["kaigara"]

EXPOSE 8080

CMD ["start", "run-loxybeat.sh"]
