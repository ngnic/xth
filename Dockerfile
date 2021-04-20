FROM golang:1.16

WORKDIR /opt/app
COPY . /opt/app

RUN go build -o github-app

CMD ["/opt/app/github-app"]