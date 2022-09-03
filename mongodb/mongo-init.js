print('Start #################################################################');
db.auth("root", "root")
db = db.getSiblingDB('organization');
db.log.insertOne({"message": "Database created."});

db.runCommand(
    {
        createUser: "one",
        pwd: "pass",
        roles: ["readWrite",],
    });

db.createCollection("event");
db.createCollection("inscription");
db.createCollection("users");

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
    ]);

db.event.createIndex({
    title: "text"
});