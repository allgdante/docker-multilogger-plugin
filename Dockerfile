FROM golang:1.16-alpine AS build
RUN apk add --no-cache gcc musl-dev
WORKDIR /go/src/github.com/allgdante/docker-multilogger-plugin
COPY . .
ENV CGO_ENABLED=0
RUN go get
RUN go build -ldflags '-extldflags "-fno-PIC -static"' -buildmode pie -tags 'osusergo netgo static_build'

FROM alpine:3.14
RUN apk add --update --no-cache ca-certificates tzdata
COPY --from=build /go/src/github.com/allgdante/docker-multilogger-plugin/docker-multilogger-plugin /bin/
WORKDIR /bin/
ENTRYPOINT ["/bin/docker-multilogger-plugin"]
