const mongo = require('mongodb').MongoClient

async function init(dbs) {
  let client = new mongo(process.env.PFF, {
    useNewUrlParser: true,
    useUnifiedTopology: true
  })

  try {
    await client.connect()
    let prices = []

    dbs.forEach((d) => {
        let db = client.db(d)
        let col = db.collection("current")

        prices[dbs.indexOf(d)] = col.find({}).toArray().map((p) => {
          return p
        })

        await db.collection("current").drop() 
        await db.collection("historical").insertMany(prices[dbs.indexOf(d)])
    })

    return true
  }

  catch (e) {
    console.error(e)
  } 
  
  finally {
    client.close()
  } 
}

init(["pff", "cef"])
    .then((res) => {
        console.log(res)
})


