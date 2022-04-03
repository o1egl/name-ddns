# name-dyndns
Dynamic DNS service for name.com

Install:

```bash
$ docker run --name name-ddns -d --restart always -e DOMAIN="test.com" -e USERNAME="user" -e TOKEN="token" o1egl/name-ddns:latest
```

| ENV      | Description                                 | Required | Default |
|----------|---------------------------------------------|----------|---------|
| DOMAIN   | Domain name                                 | yes      |         |
| HOST     | Sub domain. Keep empty for top level domain | no       |         |
| USERNAME | Your name.com usernmae                      | yes      |         |
| TOKEN    | name.com access token                       | yes      |         |
| INTERVAL | IP sync interval                            | no       | 5m      |