This project is a read-write lock for etcd. 
This requires etcd version 3.4 branch.

### To get proper etcd version
In $GOPATH, 
```
mkdir go.etcd.io; cd go.etcd.io
git clone git@github.com:etcd-io/etcd.git
git checkout release-3.4
``` 

### To run this application
- Get etcd cluster up
```
git clone git@github.com:sreeja/compose-etcd.git
cd compose-etcd
docker-compose up
```

- Run this application
```
go run main.go <logfile>
```