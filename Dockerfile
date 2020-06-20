FROM golang:buster

COPY ./main.go /main.go
COPY ./go.mod /go.mod

RUN cd / && go build -o callback

FROM debian:buster

ADD 'https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.buster_amd64.deb' "/tmp/wkhtmltox.deb"
RUN apt-get update && apt install -y "/tmp/wkhtmltox.deb" && rm -fr "/tmp/wkhtmltox.deb" 

COPY --from=0 /callback /callback
COPY ./entrypoint.sh /entrypoint.sh

USER 1337

ENTRYPOINT ["/entrypoint.sh"]