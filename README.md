## SW Engineer assignment

### Installation process
This project uses  [NATS](https://nats.io/) for communication between the server and clients. The easiest way to run it on your machine is via Docker, so
1. If you do not have Docker installed on your machine, you may follow the [instalation guide](https://docs.docker.com/get-docker/)
2. Pull NATS image:
   ```
   docker pull nats:latest
   ```
3. Run NATS in docker:
    ```
    docker run -p 4222:4222 -ti nats:latest
    ```
4. Build server and client:
    ```
    make build_all
   ```
### How to run
- After successful installation the executable binaries for client and server can be found at `./bin/` directory
- Run *server* in one terminal window and *client* in another 
- Client provides a simple command line interface for interacting with server
- The server's *"action log"* i.e. the results of client's requests can be found in `./server.log` file by default
- You may also enjoy the output in stdout ;)
- You can run as many clients as you wish :)

### TODO:
- [ ] Clients are dumb i.e. they just fire-and-forget the request. We can easily transfer to NATS Request-Response so the client's will `Get` the requested `Item(s)` and will know whether `AddItem` and `RemoveItem` operations succeeded or failed.
