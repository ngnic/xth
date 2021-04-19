FROM golang:1.16

RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash
RUN apt-get update
RUN apt-get install -y migrate

WORKDIR /opt/app
COPY . /opt/app

CMD ["/opt/app/scripts/deploy.sh"]