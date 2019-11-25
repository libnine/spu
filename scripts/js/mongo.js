const mongo = require('mongodb').MongoClient

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
    
    const gt_one = await col.find({chg_pct: {$gt: 0.99}}, 
      {ticker: 1, last: 1, chg_pct: 1, volume: 1, yield: 1}).sort({chg_pct: -1}).toArray()
    
    const lt_one = await col.find({chg_pct: {$lt: -0.99}}, 
      {ticker: 1, last: 1, chg_pct: 1, volume: 1, yield: 1}).sort({chg_pct: 1}).toArray()
    
    const rel_vol = await col.find({relative_volume: {$gt: 2}}, 
      {ticker: 1, last: 1, chg_pct: 1, volume: 1, yield: 1}).sort({relative_volume: -1}).limit(10).toArray()
    
    // const 

    return [{"gt_one_pct": gt_one, "lt_one_pct": lt_one, "rel_vol": rel_vol}]
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
    console.log(res)
  })
