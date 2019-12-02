const mongo = require('mongodb').MongoClient
const fs = require('fs-extra')

let arr = []

async function queries() {
  let client = new mongo(process.env.PFF, {
    useNewUrlParser: true,
    useUnifiedTopology: true
  })

  try {
    await client.connect()
    const db = client.db("pff")
    const col = db.collection("current")
    
    const gt_one = await col.find({chg_pct: {$gt: 0.99}, yield: {$lt: 20}, last: {$gt: 17}}).sort({chg_pct: -1}).toArray()
    const lt_one = await col.find({chg_pct: {$lt: -0.99}, last: {$gt: 17}}).sort({chg_pct: 1}).toArray()
    const rel_vol = await col.find({relative_volume: {$gt: 2}}).sort({relative_volume: -1}).toArray()
    const yld = await col.find({yield: {$gt: 6}, yield: {$lt: 20}, last: {$gt: 17}}).sort({yield: -1}).limit(10).toArray()

    return [{"gt_one_percent": gt_one, "high_yield": yld, "lt_one_percent": lt_one, "rel_vol": rel_vol}]
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
    fs.writeFileSync(`./data/dumps/dump.json`, JSON.stringify(res, null, 2))
  })
