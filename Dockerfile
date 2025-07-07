FROM --platform=$BUILDPLATFORM golang:1.24 AS builder

ARG VERSION
ARG TIMESTAMP

WORKDIR /app
COPY . /app

ARG TARGETOS TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ENV CGO_ENABLED=0
RUN go build -a -installsuffix "static" -ldflags \
    "-s -w -X github.com/openairtech/apiserver/cmd.BuildVersion=$VERSION \
           -X github.com/openairtech/apiserver/cmd.BuildTimestamp=$TIMESTAMP" \
    -o bin/openair-apiserver

FROM scratch
COPY --from=builder /app/bin/openair-apiserver /
USER 65534:65534
EXPOSE 8081
ENTRYPOINT ["/openair-apiserver"]
