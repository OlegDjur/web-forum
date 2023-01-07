FROM golang:1.18-alpine as builder

LABEL maintainer = "Oleg_Dzhur, Yernur_z - Forum"

WORKDIR /app

COPY . .

RUN   apk add build-base && go build -o main ./cmd

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app .

RUN apk add bash
EXPOSE 8000

ENTRYPOINT ["/app/main"]

# To build image:
#  docker build -f Dockerfile -t forum_image .

# To run docker container
#docker run -dp 8000:8000 --name Forum forum_image