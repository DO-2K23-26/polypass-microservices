# Polypass Microservices Interfaces

This directory contains the interface definitions for the Polypass microservices. These interfaces define the contracts between services and are used to generate code for communication between services.

## AVDL Files

AVDL (Avro IDL) files define the data structures and messages that are exchanged between services. These files are used to generate code for serializing and deserializing messages in a language-agnostic way.

### Statistics Service

The Statistics Service interfaces define the events that are consumed by the statistics service for metric calculation, as well as the metrics that are exposed by the service.

- `statistics/events.avdl`: Defines the events that the statistics service consumes, including user events, password events, and security events.
- `statistics/metrics.avdl`: Defines the metrics that the statistics service exposes, including user metrics, password metrics, and security metrics.

## Usage

### Generating Code from AVDL

To generate code from AVDL files, you can use the Avro tools:

```bash
# Install Avro tools
java -jar avro-tools-1.11.0.jar

# Generate Java code
java -jar avro-tools-1.11.0.jar idl2schemata statistics/events.avdl ./
java -jar avro-tools-1.11.0.jar compile schema com/polypass/statistics/Event.avsc com/polypass/statistics/UserLoginEvent.avsc ... ./

# Generate Go code (using gogen-avro)
gogen-avro --package statistics ./statistics com/polypass/statistics/Event.avsc com/polypass/statistics/UserLoginEvent.avsc ...
```

### Using Generated Code

The generated code can be used to serialize and deserialize messages in your services:

```go
// Go example
import (
    "github.com/polypass/polypass-microservices/interfaces/statistics"
)

// Create a new event
event := statistics.NewEvent()
event.Id = "123"
event.Type = "USER_LOGIN"
event.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
event.Source = "WEB_APP"

// Serialize the event
bytes, err := event.Serialize()
if err != nil {
    log.Fatalf("Failed to serialize event: %v", err)
}

// Deserialize the event
newEvent, err := statistics.DeserializeEvent(bytes)
if err != nil {
    log.Fatalf("Failed to deserialize event: %v", err)
}
```

## Adding New Interfaces

To add a new interface:

1. Create a new AVDL file in the appropriate directory
2. Define the data structures and messages in the AVDL file
3. Generate code from the AVDL file
4. Update this README.md file with information about the new interface
