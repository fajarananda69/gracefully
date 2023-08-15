RUN:
- Server : go run server/main.go
- Client : go run client/main.go


test gracefully
1. run server
2. run client
3. after client send request to server, kill server process (ctrl + c)
4. server will waiting until all process done or waiting 60s max time before server shutdown 