#### 1 Build docker image
```bash
# build executable bin
ego-go build hello.go

docker build -t wetee/my-app .
```

#### 2 Run docker image
```bash
docker run -p 8999:8999 --device /dev/sgx/enclave --device /dev/sgx/provision wetee/my-app
```

#### 3 open http://localhost:8999/