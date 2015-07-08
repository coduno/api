FROM scratch
MAINTAINER "Lorenz Leutgeb <lorenz.leutgeb@cod.uno>"
ADD coduno /
ENTRYPOINT ["/coduno"]
