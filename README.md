


## CompileDaemon

```
cd ~
GO111MODULE=on go get github.com/githubnemo/CompileDaemon
```
## Execution

To start the server with Docker

```
make dev
```

To stop

```
make clean
```
## Unit tests
Run unit tests on endpoints

```
make test
```
## Jenkins

This command will do exactly what Jenkins does on each push, but locally. It is recommended to run before the push.

```
make jenkins
```
## Use

Browse http://localhost:8080. The port is displayed on the console.

