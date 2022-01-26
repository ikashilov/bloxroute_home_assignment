## SW Engineer assignment

### Installation process
1. If you do not have Docker installed on your machine, you may follow the [instalation guide](https://docs.docker.com/get-docker/)
2. Pull [NATS](https://nats.io/) image:
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
- After successful installation the executable binaries for client and server can be found at: `./client/bin/client` and `./server/bin/server` respectively. And are ready for execution.  
- The server's *"action log"* i.e. the results of client's requests can be found in `./server.log` file by default  
- You may also enjoy the output in stdout ;)  
- Clients provide simple command line interface for interacting with server
- You can run as many clients as you wish :)

### TODO:
- [ ] As you may see the `internal/api` is duplicated both in client and server which is not so pretty. Can be moved to the separate package (repo)
- [ ] Clients are dumb i.e. they just fire-and-forget the request. We can easily transfer to NATS Request-Response so the client's will et the requested item(s) and will know whether `AddItem` and `RemoveItem` operations are succeeded or failed.
