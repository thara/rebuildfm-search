FROM golang:1.7

RUN apt-get update && apt-get install netcat -y

RUN curl -s https://glide.sh/get | sh

RUN wget https://github.com/Yelp/dumb-init/releases/download/v1.0.0/dumb-init_1.0.0_amd64.deb
RUN dpkg -i dumb-init_*.deb

WORKDIR /go/src/github.com/tomochikahara/rebuildfm-search

COPY run.sh /
RUN chmod +x /run.sh

ENTRYPOINT ["dumb-init", "/run.sh"]

EXPOSE 8080
