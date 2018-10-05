From golang:1.11.1

COPY . /repo
WORKDIR /repo

RUN make build

USER root

ENTRYPOINT [ "/repo/scripts/docker/entrypoint.sh" ]
