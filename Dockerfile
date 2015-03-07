FROM google/appengine-go

MAINTAINER "Lorenz Leutgeb <lorenz.leutgeb@cod.uno>"

RUN apt-get install -y -q curl build-essential git
RUN curl https://storage.googleapis.com/golang/go1.2.2.linux-amd64.tar.gz | tar xvzf - -C /goroot --strip-components=1

ENV GOROOT /goroot
ENV GOPATH /gopath
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin

# TODO(gmlewis): Remove next line once google/appengine-go image updates.
WORKDIR /app

ADD . /app

RUN echo 'machine github.com login flowlo password 04551b20222defb527351e1104c868b742db27a9' > ~/.netrc
RUN go get -d

RUN /bin/bash /app/_ah/build.sh
