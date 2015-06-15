# eaciit-mq
Memory Data Que management developed using GoLang


List commands :

1.  exit
2.  kill
3.  ping
4.  nodes
5.  gettable
6.  set(key,value)
⋅⋅* key   -> tablename|key, ex : employees|emp1
⋅⋅* value -> json format, ex : {"name":"nanda","role":"admin"}
7.  get(key)
⋅⋅* key   -> table name|key, ex : employees|emp1
8.  inc(parameter)
9.  getlog(parameter)
10. adduser(parameter)
11. deleteuser(parameter)
12. changepassword(parameter)
13. getlistusers
14. keys(nodenumber)
⋅⋅* nodenumber, ex : 0
15. info(key)
⋅⋅* key   -> tablename|key, ex : employees|emp1
16. writetodisk(key1,key2,...)
⋅⋅* key   -> tablename|key, ex : employees|emp1
17. readfromdisk(key1,key2,...)
⋅⋅* key   -> tablename|key, ex : employees|emp1
