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
## Kaigara ]>

VOLUME ["${APP_HOME}"]

COPY loxybeat /usr/local/bin/
COPY run-loxybeat.sh /usr/local/bin/
COPY loxybeat.yml "${APP_HOME}"
COPY loxybeat.template.json "${APP_HOME}"
COPY loxybeat.template-es2x.json "${APP_HOME}"

ENTRYPOINT ["kaigara"]

CMD ["start", "run-loxybeat.sh"]
#CMD ["start", "bash"]
