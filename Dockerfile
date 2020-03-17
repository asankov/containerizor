FROM golang:1.13

RUN mkdir -p /app

WORKDIR /app

ADD . /app

# download all dependencies to local cache
RUN go mod download

RUN apt-get update

RUN apt-get install apt-transport-https ca-certificates curl

RUN curl -fsSl https://download.docker.com/linux/debian/gpg | apt-get key add -


RUN apt-get install docker-ce docker-ce-cli containerd.io

# build the executable binary
RUN go build -o app/cmd/web/app github.com/asankov/containerizor/cmd/web

CMD ./app/cmd/web/app -port=4000 -db_host=$DB_HOST -db_port=$DB_PORT -db_user=$DB_USER -db_name=$DB_NAME -db_pass=$DB_PASS