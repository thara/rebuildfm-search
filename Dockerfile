FROM golang:1.7
ARG command

RUN apt-get update && apt-get install netcat -y

RUN curl -s https://glide.sh/get | sh

RUN wget https://github.com/Yelp/dumb-init/releases/download/v1.0.0/dumb-init_1.0.0_amd64.deb
RUN dpkg -i dumb-init_*.deb

COPY run.sh /
RUN chmod +x /run.sh
RUN mkdir -p /go/src/github.com/thara/rebuildfm-search
WORKDIR /go/src/github.com/thara/rebuildfm-search

COPY glide.lock /go/src/github.com/thara/rebuildfm-search
COPY glide.yaml /go/src/github.com/thara/rebuildfm-search
COPY static /go/src/github.com/thara/rebuildfm-search/static
COPY templates /go/src/github.com/thara/rebuildfm-search/templates
COPY rebuildfm /go/src/github.com/thara/rebuildfm-search/rebuildfm
COPY main.go /go/src/github.com/thara/rebuildfm-search

ENV COMMAND $command
ENTRYPOINT ["dumb-init", "/run.sh"]
