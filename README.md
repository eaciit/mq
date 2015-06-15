# eaciit-mq
Memory Data Que management developed using GoLang


List commands :

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


Format key,value,nodenumber :

1.  key   -> tablename|key, ex : employees|emp1
2.  value -> json format, ex : {"name":"nanda","role":"admin"}
3.  nodenumber -> ex : 0
