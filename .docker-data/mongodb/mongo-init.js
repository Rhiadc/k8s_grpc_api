db.createUser(
        {
            user: "teste",
            pwd: "testepass",
            roles: [
                {
                    role: "readWrite",
                    db: "testedb"
                }
            ]
        }
);

db.myCollectionName.insert({ "address": { "city": "Paris", "zip": "123" }, "name": "Mike", "phone": "1234" });