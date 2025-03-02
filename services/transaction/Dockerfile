FROM golang:1.24-alpine AS base
WORKDIR /app
ENV NAME=transaction

ENV GOPATH=/go CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN apk update && apk --no-cache add git ca-certificates && update-ca-certificates

COPY . .

RUN go mod tidy && go build  -C ./cmd/ -o $NAME && chmod 777 ./cmd/${NAME}


FROM alpine:latest AS final 
RUN apk update && apk add ca-certificates
ENV NAME=transaction
CMD $NAME
COPY --from=base /app/cmd/$NAME /bin/