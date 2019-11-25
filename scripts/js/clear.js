const mongo = require('mongodb').MongoClient

let arr = []

async function init() {
  let client = new mongo(process.env.PFF, {
    useNewUrlParser: true,
    useUnifiedTopology: true
  })

  try {
    await client.connect()
    const db = client.db("pff")
    const col = db.collection("symbols")
    col.find({}).sort({"symbol": 1}).toArray((e, d) => {
        if (e) throw e
        arr = d
    })
    db.collection("historical").insertMany(arr)
    return db.collection("current").drop((e, r) => {
        if (e) throw e
        return r
    })
  } 
  
  catch (e) {
    console.error(e)
  } 
  
  finally {
    client.close()
  } 
}



init()
    .then((ok) => {
        console.log(ok)
    })