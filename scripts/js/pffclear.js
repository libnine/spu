use pff
db.current.find({}).forEach(function(d) {
    db.historical.insert(d)
})
db.current.drop()