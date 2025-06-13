# Confluent Schemas

This directory contains the schemas used by the Confluent platform for Kafka. Each kafka topic has its own schema, which is used to validate the data being sent to and received from the topic. The schemas are defined in Avro format and they are stored in the Schema Registry. The Schema Registry is a service that provides a RESTful interface for managing Avro schemas and allows for schema evolution.

## Folder structure
Each schemas has its own folder, which contains the following files:
- `configmap.yaml`: This file contains the schema definition in Avro format.
- `schema.yaml`: This file configures the schema in the Schema Registry (deployed in [the previous readme](../../README.md)).

## How to use
To use the schemas, you need to apply the configmap and schema files to your Kubernetes cluster. You can do this by running the following command:
```sh
kubectl apply -f configmap.yaml
kubectl apply -f schema.yaml
```

## How to check the schemas
You can check the schemas in the Schema Registry by portforwarding the service to your localhost:
```sh
kubectl port-forward -n confluent svc/schema-registry 8081:8081
```
Then you can access the Schema Registry at `http://localhost:8081` and check the schemas with diffrents api endpoints.

### Schemas
- List all schemas: `http://localhost:8081/schemas`
- Get a specific schema: `http://localhost:8081/schemas/ids/{schema_id}`

### Subjects
- List all subjects: `http://localhost:8081/subjects`
- Get a specific subject: `http://localhost:8081/subjects/{subject_name}`
- Get the number of versions of a subject: `http://localhost:8081/subjects/{subject_name}/versions`
- Get the `n` version of a subject: `http://localhost:8081/subjects/{subject_name}/versions/{version}`
- Get the latest version of a subject: `http://localhost:8081/subjects/{subject_name}/versions/latest`