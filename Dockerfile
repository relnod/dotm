From golang:1.11 as build 

COPY . /build
WORKDIR /build

RUN make build

FROM scratch

COPY --from=build /build/dotm /dotm

ENTRYPOINT [ "/dotm" ]
