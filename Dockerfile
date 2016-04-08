FROM kiasaki/alpine-golang

ADD settings.yml $GOPATH/bin/settings.yml
RUN go get github.com/herald-it/goncord

EXPOSE 10001

CMD ["goncord"]