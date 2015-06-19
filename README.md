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

##API

Protocol | URI | Action
--- | --- | ---
POST | `/api/gettoken/username={username}&password={password}` | Return token and valid time with given username and password
GET | `/api/checktoken/token={token}` |  Return token and time with given token
GET | `/api/get/token={token}&key={key}` | Return item with given key
POST | `/api/put/token={token}&key={key}` | Set item with given key + form value and return item info

###Get

####Example

Given Parameter

`http://localhost:8090/api/get/token=I14jVFA5UFw6LBlRWlswBGA-Lwc7DxhbO1VJNTshKRU=&key=eat`

Return Data

```
{
  "message": "",
  "data": {
    "Created": "0001-01-01T00:00:00Z",
    "Duration": 0,
    "Expiry": 0,
    "Key": "public|eat",
    "Owner": "public",
    "Permission": "666",
    "Size": 8,
    "LastAccess": "2015-06-19T10:26:05.78013317+07:00",
    "Table": "",
    "Value": "sushi"
  },
  "success": true
}
```

###Put

####Example

Given Parameter

`http://localhost:8090/api/put/token=I14jVFA5UFw6LBlRWlswBGA-Lwc7DxhbO1VJNTshKRU=&key=eat`

With FormData

`"value" = "sushi"`

Return Data

```
{
  "message": "",
  "data": {
    "Node": 0,
    "Owner": "public",
    "Size": 8,
    "Valid": 0
  },
  "success": true
}
```

###Get Token

####Example

Given Parameter

`http://localhost:8090/api/gettoken/username=set&password=rahasia`

Return Data

```
{
    "data": {
        "Token": "DDVFEhQhB1ZOJ0cbGQdbCRseOz4IXU4gOUYUGkIkNSU=",
        "Valid": "2015-06-19T11:21:36.566005061+07:00"
    },
    "message": "",
    "success": true
}
```

###Check Token

####Example

Given Parameter

`http://localhost:8090/api/checktoken/token=DDVFEhQhB1ZOJ0cbGQdbCRseOz4IXU4gOUYUGkIkNSU=`

Return Data

```
{
    "data": {
        "Token": "DDVFEhQhB1ZOJ0cbGQdbCRseOz4IXU4gOUYUGkIkNSU=",
        "Valid": "2015-06-19T11:25:52.788241816+07:00"
    },
    "message": "",
    "success": true
}
```

