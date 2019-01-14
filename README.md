# porth
A simple secure tcp port redirector.


Build
=====
### Certificates

The certificates and keys have been generated for test/development purposes only. Do not use these files in production deployments! You can regenerate them anytime by running the following command.

    cd ./tools
    ./keys.sh

### Binaries
    # build server
    cd ./server
    go build
        
    # build client
    cd ./client
    go build
     
    # build operator
    cd ./operator
    go build

Use
=====


Start the server. Server default port is 2671.
    
    ./server/server
    
Add a redirection for port 7000 on client host.

    ./client/client -server=localhost:2671 -local=:7000 -id=testid

Redirect local port 7001 on operator host using the same -id.

    ./operator/operator -server=localhost:2671 -local=:7001 -id=testid

Now the client port 7000 is accessible from operator host on port 7001.

Server status is available on port 2672.

    curl -i http://localhost:2672/clients
    curl -i http://localhost:2672/server_status

