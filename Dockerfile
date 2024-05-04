FROM golang:1.22 as build-stage

WORKDIR /app

# TODO: Optimize COPY to only copy code. Note: not urgent on ephemeral stage
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /mdroid-monolith ./cmd/monolith/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o /mdroid-yt-uploader ./cmd/yt-upload-v2/main.go && \
    CGO_ENABLED=0 GOOS=linux go build -o /mdroid-discord-bot ./cmd/discord_bot/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

ENV TZ=America/Montreal

WORKDIR /

COPY --from=build-stage /mdroid-monolith /mdroid-monolith
COPY --from=build-stage /mdroid-yt-uploader /mdroid-yt-uploader
COPY --from=build-stage /mdroid-discord-bot /mdroid-discord-bot

COPY --from=build-stage /app/assets /assets

EXPOSE 3333

USER nonroot:nonroot

CMD ["/mdroid-monolith"]