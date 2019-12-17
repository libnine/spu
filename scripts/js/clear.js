const mongo = require('mongodb').MongoClient

async function init(mdb) {
  let client = new mongo(process.env.PFF, {
    useNewUrlParser: true,
    useUnifiedTopology: true
  })

  try {
    await client.connect()
    const db = client.db(mdb)
    const col = db.collection("current")
    let colarray = await col.find({}).toArray()
    await db.collection("historical").insertMany(arr)
    return db.collection("current").drop()
  } 
  catch (e) {
    console.error(e)
  } 
  
  finally {
    client.close()
  } 
}

init("pff")
  .then((res) => {
    console.log(res)
})

init("cef")
  .then((res) => {
    console.log(res)
})


