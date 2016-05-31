FROM alpine:latest

WORKDIR /app

ADD settings.yml settings.yml
ADD ./build/goncord .

EXPOSE 8000

CMD ["./goncord"]
