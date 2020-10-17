FROM golang:1.15 AS build
RUN mkdir -p /home/main
WORKDIR /home/main

# Deps
ENV GO11MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

# Build
RUN cd cmd/server && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o buffService && \
    cp buffService /

# Runner
FROM alpine
RUN apk update && \
    apk upgrade && \
    apk add ca-certificates && \
    update-ca-certificates && \
    apk add --update tzdata && \
    apk add curl && \
    rm -rf /var/cache/apk/*

COPY --from=build /buffService /home/
ENV TZ=Europe/London
ENV PORT=80

# Entrypoint
WORKDIR /home
RUN echo "#!/bin/bash" > ./entrypoint.sh
RUN echo "./buffService" >> ./entrypoint.sh
RUN chmod +x ./entrypoint.sh

# EntryPoint
ENTRYPOINT ["sh", "./entrypoint.sh"]

HEALTHCHECK --interval=5s --timeout=2s --retries=10 CMD curl --silent --fail localhost/probe || exit 1

EXPOSE 80
