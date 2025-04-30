# Authorization service
---

## Launch the application

Populate the env.local.sh file then :
```sh
source env.local.sh
```

Start the external services:
```sh
docker-compose up -d
```

Start the application:
```sh
go run main.go
```

## Purpose

This service will hold the logic to do data aggregation form other services.
We will use [Authzed](https://authzed.com/) to manage authorization.

## Main feature

The main feature of this service is to provide a centralized authorization logic for all the microservices in the Polypass ecosystem.
We will listen to kafka topics to aggregate the authorization datas in Authzed database.
Then we will expose the GRPC api of Authzed so that it can be called by other services.

## Authorization Logic

Need authorization logic on
- folder
- tag
- credential

Role are hold at folder level. Either a user is a folder user either it a folder viewer or  a folder viewer.
A folder admin will inherit all the folder permission that as a folder viewer.

Basically the folder viewer will only be able to perform read operation on:
- folder
- tag
- credential
- sub-folders
- credentials in sub-folders
- tags in sub-folders

The folder admin admin will have the right to:
- create sub-folders
- update folder and sub-folders
- create & update tags in the folder and sub-folders
- create & update  credentials in the folder and sub-folders

## Some reliationship writing help

For example, if you want to add some credential to a folder you will have to write this reliationship:

```
Resource Type: credentials
Resource Id: "ascsd1234567890"
Subject Type: folder
Subject Id: "fo1234567890"
```

## Open the playground Authzed

```sh
docker run -it -p 3000:3000 ghcr.io/authzed/spicedb-playground:latest
```

Open the playground Authzed

## Discussion

It could be interesting to let the check api of Authzed directly exposed in order to not re-implement in this service.
