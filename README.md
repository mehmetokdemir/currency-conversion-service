## Currency Conversion Service

Build locally;
```shell
make build
```

Build as docker image;
```shell
docker build -t ${username}/${registry}:${tag} .
```

Up to database with docker compose;
````shell
docker-compose up -d
````

Run locally with go command;
````shell
make rwgocommand
````

Run locally with executable table;
````shell
make run
````

Run with docker;
````shell
docker-compose build 
docker-compose up          
````

Generate api docs;
````shell
make generate-docs
````

Swagger api documentation;
`*/swagger/index.html`


Start golangci lint run 
````shell
make lint
````