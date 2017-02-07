# rebuildfm-search
WIP


## Deploy to DigitalOcean

```
$ docker-machine create --driver digitalocean \
    --digitalocean-access-token $DIGITALOCEAN_TOKEN \
    --digitalocean-region sgp1 \
    --digitalocean-size 4gb \
    --digitalocean-image ubuntu-16-04-x64 \
    docker-prod
$ docker-machine scp -r public docker-prod:/root/
$ docker-machine scp -r rebuildfm docker-prod:/root/
$ docker-machine scp main.go docker-prod:/root/
$ docker-machine scp glide.lock docker-prod:/root/
$ docker-machine scp glide.yml docker-prod:/root/
$ docker-compose $(docker-machine config docker-prod) build
$ docker-compose $(docker-machine config docker-prod) up -d
```
