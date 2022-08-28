## _$mongosh_

```sh
use organization
```
```sh
db.createCollection("event")
```
```sh
db.createCollection("inscription")
```
```sh
db.createCollection("users")
```
```sh
db.users.insertMany([
        {
         "username": "usuarioNormal",
         "email": "usuario@test.com",
         "password": "123456",
         "type": "user",
         "token": "",
        },
        {
         "username": "usuarioAdministrador",
         "email": "admin@test.com",
         "password": "123456",
         "type": "admin",
         "token": "",
        },
        ])
```
```sh
    db.getSiblingDB("organization").runCommand(
    {
        createUser: "one",
        pwd: "pass",
        roles: ["readWrite",],
    })
```
