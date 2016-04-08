FROM alpine:latest

WORKDIR /app

ADD settings.yml settings.yml
ADD ./build/goncord .

EXPOSE 10001

CMD ["./goncord"]