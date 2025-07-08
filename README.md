# Go Example #

Create a program to explore Go' language.

## Create an L4 Load Balancer

### Requirements

- An input connection can manage multiple output instances
- Can handle unhealthy instances and watch them
- Can manage at least TCP connections
- One port can serve multiple requests
- It manages requests and responses
- Can manage multiple connections in different ports

### Concetps

#### Client
Who open a socket to send and recive data through the connections

#### Connection
Is the served connection for the clients

#### Service
The service who recive the client request and process it to give a response

##### Instance
Each "real" instance of a service with an address


### Process

1. Service starts
2. Read the configuration and generate the current objects
3. Start the sockets and redirect all incomming traffic to the service's instance
