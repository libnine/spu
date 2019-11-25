const mongo = require('mongodb').MongoClient
const fs = require('fs-extra')
const path = require('path')

let arr = []

async function queries() {
  let client = new mongo(process.env.PFF, {
    useNewUrlParser: true,
    useUnifiedTopology: true
  })

  try {
    await client.connect()
    const db = client.db("pff")
    const col = db.collection("bu")
    
    const gt_one = await col.find({chg_pct: {$gt: 0.99}}).sort({chg_pct: -1}).toArray()
    const lt_one = await col.find({chg_pct: {$lt: -0.99}}).sort({chg_pct: 1}).toArray()
    const rel_vol = await col.find({relative_volume: {$gt: 2}}).sort({relative_volume: -1}).limit(10).toArray()
    
    return [{"gt_one_percent": gt_one, "lt_one_percent": lt_one, "rel_vol": rel_vol}]
  }

  catch (e) {
    console.error(e)
  }

  finally {
    client.close()
  }
}

queries().
  then((res) => {
    fs.emptyDirSync("./data/dumps/")

    let n = Math.floor(Math.random() * 100000000000)
    fs.writeFileSync(`./data/dumps/${n}_mongo.json`, JSON.stringify(res, null, 2))
  })
