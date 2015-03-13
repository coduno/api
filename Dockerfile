FROM fedora:21

MAINTAINER "Lorenz Leutgeb <lorenz.leutgeb@cod.uno>"

RUN mkdir -v /go
ENV GOPATH /go

RUN yum -y install curl tar git

RUN curl https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz | tar --exclude='go/doc' --exclude='go/blog' -xvzf - -C /usr/local
ENV PATH $PATH:/usr/local/go/bin

WORKDIR /app
ADD . /app

RUN echo 'machine github.com login flowlo password 04551b20222defb527351e1104c868b742db27a9' > ~/.netrc
RUN go get -d
RUN go build -work -x -v -o coduno

ENTRYPOINT ["/app/coduno"]
