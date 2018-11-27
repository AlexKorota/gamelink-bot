FROM alpine

WORKDIR /app

COPY gamelink-bot user.html ./

RUN apk update && apk add --no-cache ca-certificates

ENTRYPOINT [ "./gamelink-bot" ]