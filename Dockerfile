FROM golang:1.5
MAINTAINER Eagle Chen <chygr1234@gmail.com>

COPY metrics.go /tmp/
COPY Godeps /tmp/

RUN \
  go get github.com/tools/godep && \
  cd /tmp && \
  godep go install

ENTRYPOINT ["/docker-entrypoint.sh"]

EXPOSE 12801
