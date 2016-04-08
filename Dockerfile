FROM ubuntu:latest

WORKDIR /auth-service

ADD auth-service /auth-service
ADD settings.yml /auth-service

EXPOSE 10001

CMD ["./auth-service"]