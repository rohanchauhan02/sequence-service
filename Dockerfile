FROM public.ecr.aws/docker/library/golang:1.25.1-alpine3.22 AS build_base

RUN apk update && apk upgrade && \
    apk --no-cache --update add make gcc g++ libc-dev

WORKDIR /go/src/github.com/rohanchauhan02/sequence-service
ENV GO111MODULE=on
ENV GODEBUG="madvdontneed=1"

COPY go.mod .
COPY go.sum .
RUN go mod download

FROM build_base AS server_builder

WORKDIR /go/src/github.com/rohanchauhan02/sequence-service
COPY . .

RUN GOOS=linux GOARCH=amd64 make build

FROM alpine:latest

RUN apk update && apk upgrade

WORKDIR /sequence-service/app
EXPOSE 8880

# Copy binary and configs
COPY --from=server_builder /go/src/github.com/rohanchauhan02/sequence-service/ .
COPY --from=server_builder /go/src/github.com/rohanchauhan02/sequence-service/configs ./configs/
COPY --from=server_builder /go/src/github.com/rohanchauhan02/sequence-service/files ./files/

CMD /sequence-service/app/engine
