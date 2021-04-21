FROM golang:1.16

RUN apt update
RUN apt install -y postgresql

WORKDIR /opt/app
COPY . /opt/app

RUN go build -o github-app

CMD ["/opt/app/github-app"]