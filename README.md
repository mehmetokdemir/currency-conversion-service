## Currency Conversion Service

Build locally;
```shell
make build
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
make run-with-docker       
````

Generate api docs;
````shell
make generate-docs
````

Swagger api documentation;
````shell
*/swagger/index.html
````

Start golangci lint run 
````shell
make lint
````