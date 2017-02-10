# Rebuild.fm Search

Rebuild.fm Search is a search service for [Rebuild.fm](https://rebuild.fm)'s episodes.


## Build JS

```
$ elm-make elm/Search.elm --output static/elm.js
```


## Deploy to DigitalOcean

```
$ docker-machine create --driver digitalocean \
    --digitalocean-access-token $DIGITALOCEAN_TOKEN \
    --digitalocean-region sgp1 \
    --digitalocean-size 4gb \
    --digitalocean-image ubuntu-16-04-x64 \
    docker-prod
$ docker-machine scp -r static docker-prod:/root/
$ docker-machine scp -r templates docker-prod:/root/
$ docker-machine scp -r rebuildfm docker-prod:/root/
$ docker-machine scp main.go docker-prod:/root/
$ docker-machine scp glide.lock docker-prod:/root/
$ docker-machine scp glide.yml docker-prod:/root/
$ export SITE_URL=http://rebuildfm-search.thara.jp
$ docker-compose $(docker-machine config docker-prod) build
$ docker-compose $(docker-machine config docker-prod) up -d
```
