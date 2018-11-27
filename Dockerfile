FROM alpine

WORKDIR /app

COPY gamelink-bot ./

RUN apk update && apk add --no-cache ca-certificates

ENTRYPOINT [ "./gamelink-bot" ]