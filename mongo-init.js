db = db.GetSiblingDB('admin')

db.auth("root", "rootpassword")

db = db.GetSiblingDB('smolathon')
db.CreateUser({
    'user': "user",
    'pwd': "user",
    'roles': [{
        'role': 'dbOwner',
        'db': 'smolathon'
    }]
});
db.CreateCollection('init');