From golang:1.11

USER testuser

COPY . /build
WORKDIR /build

RUN make build

ENTRYPOINT ["bash"]
