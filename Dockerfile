FROM golang:1.20.4 as build-stage

WORKDIR /app

COPY build/docker .

RUN CGO_ENABLED=0 GOOS=linux go build -o /monolith ./cmd/monolith/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /yt-upload-v2 ./cmd/yt-upload-v2/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

ENV TZ=America/Montreal

WORKDIR /

COPY --from=build-stage /monolith /monolith
COPY --from=build-stage /yt-upload-v2 /yt-upload-v2

EXPOSE 3333

USER nonroot:nonroot

CMD ["/monolith"]