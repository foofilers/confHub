FROM golang:1.8.1
MAINTAINER Igor Maculan <n3wtron@gmail.com>

COPY ../../ /go/src/github.com/foofilers/confHub/

WORKDIR $GOPATH/src/github.com/foofilers/confHub/

RUN go-wrapper download
RUN go-wrapper install

CMD ["go-wrapper", "run"]