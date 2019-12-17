const mongo = require('mongodb').MongoClient

prices = []

async function init(dbs) {
  let client = new mongo(process.env.PFF, {
    useNewUrlParser: true,
    useUnifiedTopology: true
  })

  try {
    await client.connect()

    dbs.forEach((d) => {
        let db = client.db(d)
        let col = db.collection("current")

        prices[dbs.indexOf(d)] = col.find({}).toArray().then((p) => {
            p.forEach((price) => {
                return price
            })
            db.collection("current").drop()
        })
        db.collection("historical").insertMany(prices[dbs.indexOf(d)])
    })

    return prices.length
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


