# url-shortener

### Getting Start

## Redis
This application use Redis as its primary datastore. Either specify a redis cluster in config to run one on your local with docker.

You can start a Redis instance on your local by running this command.

```
docker run --name redis -p 6379:6379 -d redis
```

## Config File
This application has a sample config file in its root. Move this file to an appropriate directory and specific that directory in the config.go file.

## Install Dependencies
You can install all needed dependencies by running this command in the root directory.

```
go mod download
```

## Run application
Finally, you can run this appliation by running this command in the root directory.

```
go run .
```