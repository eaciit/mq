# eaciit-mq
Memory Data Que management developed using GoLang

##Running Server / Node

###Running Node as Master

```
go run mqd.go
```

Node will automatically run as server on localhost:7890

###Running Node as Slave

```
go run mqd.go -master 127.0.0.1:7890 -port 7891
```

Node will automatically run as slave with master on localhost:7890
Note : slave port must different with master port

###Running Node as Mirror

```
go run mqd.go -master 127.0.0.1:7890 -port 7892 -mirror
```

Adding ```-mirror``` will run node as mirror
Node will automatically run as mirror with master on localhost:7890
Note : mirror port must different with master port

## Running Web Monitor
```
go run mqmonitor.go
```
Automatically run web server on localhost:1234 and connecting to master node on 127.0.0.1:7890

## Running Client
```
go run mqclient.go
```
Automatically run client and connecting to RPC Server.

###List client commands
Command | Action
--- | ---
`exit` | exit from mqclient
`kill` | kill all nodes and client (webmonitor will `still` running)
`ping` | show recent status of all nodes (including mirror)
`nodes` | -
`gettable` | -
`set(key,value)` | set value for given key or adding new value with given key if not existed
`get(key)` | get value for given key
`inc(parameter)` | -
`getlog(parameter)` | -
`addUser(username,password,role)` | add new user with given parameters *(Only admin can run this command)*
`updateUser(username,password,role)` | update registed user's password and role form given username *(Only admin can run this command)*
`deleteUser(username1,username2,...)` | delete user with given usernames *(Only admin can run this command)*
`changepassword(newPassword)` | change current password to new password
`getlistusers` | show all username and encrypted password registed
`keys(nodenumber)` | show all available keys on given node
`info(key)` | show detailed info for given key
`writetodisk` | write all available items to node's disk
`writetodisk(key1,key2,...)` | write items for given keys to node's disk
`readfromdisk` | read all available data on all node's disk
`readfromdisk(key1,key2,...)` | *Not yet implemented*

###Format key,value,nodenumber :
```
1.  key   -> tablename|key, ex : employees|emp1
2.  value -> json format, ex : {"name":"eaciit","role":"admin"}
3.  nodenumber -> ex : 0
```
