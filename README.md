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

##List client commands :

```
1.  exit
2.  kill
3.  ping
4.  nodes
5.  gettable
6.  set(key,value)
7.  get(key)
8.  inc(parameter)
9.  getlog(parameter)
10. adduser(parameter)
11. deleteuser(parameter)
12. changepassword(parameter)
13. getlistusers
14. keys(nodenumber)
15. info(key)
16. writetodisk(key1,key2,...)
17. readfromdisk(key1,key2,...)
```
###Format key,value,nodenumber :
```
1.  key   -> tablename|key, ex : employees|emp1
2.  value -> json format, ex : {"name":"nanda","role":"admin"}
3.  nodenumber -> ex : 0
```

##User Management Command

###Adding new user

``` 
addUser(username,password,role) 
```
example: 
``` 
addUser(eaciit,master,admin) 
```
this command will add new user with username: "eaciit", password: "master", role: "admin"
Note : only admin can add new user

###Change Password of Current User

``` 
changePassword(newPassword)
```
example: 
``` 
changePassword(secret)
```
this command will change current user password to secret

###Delete existing user

``` 
deleteUser(username1,username2,username3,...)
```
example: 
``` 
deleteUser(eaciit1,eaciit2) 
```
will delete user with username "eaciit1" and "eaciit2"
Note : only admin can delete existing user

###Updating new user

``` 
updateUser(username,password,role) 
```
example: 
``` 
updateUser(eaciit,secret,admin) 
```
will update password and role of existing user with username  "eaciit" to password: "secret" and role: "admin"
Note : only admin can update existing user


