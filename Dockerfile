FROM golang:tip-alpine as firstbeego

WORKDIR /app
COPY . .

RUN go build -tags netgo -o firstbeego .

FROM alpine:latest

RUN apk update && apk add --no-cache git

COPY --from=firstbeego /app/firstbeego /firstbeego

COPY --from=firstbeego /app/conf /conf

WORKDIR /

CMD ["/firstbeego"]