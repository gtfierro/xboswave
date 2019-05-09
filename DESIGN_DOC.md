# XBOS(2) Design Document

How to organize:

## Services + Interfaces

**Services** are computational processes that provide functionality over the network.
Examples include:

- databases (providing data storage and retrieval)
- models (providing predictions)
- controllers (providing control decisions)
- drivers (providing getter/setter interfaces to cyberphysical resources)

Services expose one or more **interfaces**. Interfaces define **signals** and **slots**. We can look at this from an object-oriented programming perspective. Services are objects, interfaces are like Go interfaces or Java abstract classes. Signals are analagous to "getter" methods, and slots are analogous to "setter" methods.

### IDL

All messages in XBOS use the protobuf IDL.
The `proto` files for services in XBOS are distributed as part of the release. At some point, they may be persisted as messages on a well-known namespace so that clients can easily find the latest files and regenerate their bindings.
Protobuf messages are designed to grow naturally over time, so clients should be able to interact with messages adhering to more recent protobuf definitions without breaking.

Protobuf messages need to account for device/service/interface heterogeneity. We should be able to compose messages together using the graph-based description.

### Liveness

Services periodically announce their liveness.

#### Questions

Done using a persisted message on WAVEMQ? Or just published? What if this is not persisted?

- harder to list services by doing the query on WAVEMQ
- but should we just be using the Brick model for this anyway?
- but then how would the Brick model get generated?

### URI Structure

We introduce a transport agnostic structure for describing + identifying services:

A service resource is described by:

- a `location` (IP address or WAVEMQ URI)
- a `namespace` (root of authority)
- a `service` (pkg name + service name)
- `instance ID`
- `signal`/`slot`
- a `method` name

This incorporates elements of both GRPC URL structure and BW2-style service+interface.

| XBOS | BW2 | GRPC |
| --- | --- | --- |
| `location` -- possibly remove this? | namespace + base URI | n/a |
| `namespace` | namespace | n/a |
| `service` | `s.servicename` | pkg name + service name |
| `instance ID` | `instance id` | n/a |
| `signal` / `slot` | `signal` / `slot` | call method / respond |
| `method` | `signal`/`slot` name | method name |

Users will call + refer to services, drivers etc using this structure.
For permissions, this structure will "compile down" into the following WAVE URI structure:

```
Namespace: namespace
Resource: xbos/<location>/<service>/<instance ID>/<signal/slot>/<method>
```

This should make it easier to develop "templates" for more easily determining what permissions are needed to interact with a service/device/etc.
When client libraries are interacting with these structures, they can determine if they need to dial the IP address to interact over GRPC or if they need to subscribe/publish on WAVEMQ.

### Models of Communication

There are two models of communication that we make use of in XBOS: publish-subscribe and request-response.
Publish-subscribe is the dominant model for reporting state and simple commands:

Publish-subscribe:
- publishing:
    - sensor reading
    - device state
    - service stats:
        - database query latency histogram
        - prediction error for models
        - service uptime
    - heartbeats?
    - contextual information
        - Brick device + point relationships
        - service descriptions (controller ties switch to light)
- subscribing:
    - device actuations:
        - using messages to account for cases where multiple fields must be set simultaneously
        - could use arrays of actuation messages to convey control sequences
    - ingester:
        - archive fields in subscribed messages
    - alarm service:
        - OOB alert when conditions are broken

Request-response:
- invoking services:
    - database query
    - model prediction

To make sense of when to use either implementation, we can say that *most* services will communicate
using WAVEMQ.  Some will have the base setter/getter, some will use RPC, and a few "well-known" services
may choose to use to provide GRPC interface w/o WAVEMQ (but using WAVE permissions).

Whatever method we choose for storing the URIs/locations of services+drivers needs to distinguish between
which of these is to be used:

- WAVEMQ getter+setter:
    - content: needs protobuf message definitions
    - location: needs namespace + URIs
- WAVEMQ RPC:
    - content: needs protobuf message definitions + service+rpc call
    - location: needs namespace + URIs
- WAVE RPC:
    - content: needs protobuf message definitions + service+rpc call
    - permissions: needs namespace + URI
    - location: URL

- two implementations for WAVEMQ RPC:
    - use WAVE to authorize/authenticate GRPC:
        - advantages:
            - server side does not need any changes
            - get benefits of WAVE delegation model
        - disadvantages:
            - if GRPC service is behind firewall/LAN, we cannot invoke it directly
    - use GRPC over WAVEMQ
        - advantages:
            - consistent communication plane
            - full location transparency
            - same firewall behavior
        - disadvantages:
            - not as performant, likely a bit buggy at first
            - requires altering servers
                - unless we do it as a proxy

## Service Discovery

- services + drivers publish little descriptions of themselves
    - these get persisted on the message bus?
        - form of multicast e.g. upnp
        - do we want to periodically announce these? heartbeats?
        - this is in contrast to a database that holds all known services
    - these describe how to access the service:
        - WAVEMQ URIs or IP address (for GRPC service)
        - GRPC/protobuf links, files
        - other information?
    - would we need to negotiate anything between the client, or is this
      just a consumable service description?

## Composeable Service Definitions

Need to unify the graph description of services + drivers + their context with the actual messages they send.

Brick model of a thermostat:

```
:tstat1     a               brick:Thermostat
:tstat1     bf:hasPoint     :tempsen1
:tstat1     bf:hasPoint     :hsp1
:tstat1     bf:hasPoint     :csp1
:tempsen1   a               brick:Temperature_Sensor
:hsp1       a               brick:Heating_Temperature_Setpoint
:csp1       a               brick:Cooling_Temperature_Setpoint
```

Each Brick superclass has a protobuf message definition. Subclasses inherit the message definitions of their parent classes.

```proto
message Temperature_Sensor {
    int64   time = 1;
    Double  value = 2;
}

# superclass of heating/cooling sp
message Temperature_Setpoint {
    int64   time = 1;
    Double  value = 2;
}

# can "compile" the more specific classes that aren't explicitly defined
# from their super classes
## DERIVED
# message Heating_Temperature_Setpoint {
#    int64   time = 1;
#    Double  value = 2;
# }
# message Cooling_Temperature_Setpoint {
#    int64   time = 1;
#    Double  value = 2;
# }

message Thermostat {
    int64 time = 1;
    Temperature_Sensor temperature = 2;
    Heating_Temperature_Setpoint hsp = 3;
    Cooling_Temperature_Setpoint csp = 4;
}
```

If we have the constraint that each "device" only has one field of each Brick class, then we can
unambiguously go between definitions.

#### Questions

- what are the use cases for going between the two representations?
    - execute a brick query to get e.g. all temperature sensors. Want to subscribe to them, but they exist as parts of thermostats, hamiltons and other temperature sensors. This means that the messages containing the fields we actually care about are of all different types. How do we rectify this?
- can we make the assumption that each message will only have at most one field of each Brick type?

## `xbosd`

Single binary that provides most of the functionality needed for the local network:

- embedded `wave` API
- embedded `wavemq` API
- embedded XBOS-specific abstraction API?
    - one that consumes the service specification above?

Additional methods on the binary:

- `xbos cmd`: runs a script with the proper environment variables set
- `xbos daemon`: manages persistent services
- `xbos update`: updates services, etc


## Outline

- services + interfaces for computational processes
    - semantics:
        - drivers: set state / get state
        - services: request / response
    - mechanisms:
        - determining liveness (heartbeats)
        - determining URI structure
        - reporting status
        - discovery
            - WAVEMQ query/subscribe
            - Brick model
        - context management:
            - how to announce semantic description of service/driver
                - this is NOT the physical context. Just description of attributes
            - how to download semantic context of where driver "is"
        - communication:
            - request-response driven: (use GRPC or GRPC / WAVEMQ)
                - single recipient
            - event-driven: use WAVEMQ pubsub
                - multiple recipient
                - how to rectify actuation messages? Aren't they just req/resp?
                  The results of the requests are multiple-recipient
- semantic context
    - directed graph-based:
        - rooted in Brick ontology
        - expand to incorporate service relationships
    - how to form/populate/update
    - how to store
        - who stores it
    - role in permissions management
    - how to apply permissions
- templates/descriptions of common services + drivers and how they fit into this
    - thermostat driver
    - light switch driver
    - light (bulb) driver
    - ingester
    - timeseries database
    - brick model database
