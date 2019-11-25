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
    return await col.drop((e, r) => {
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


const axios = require('axios')
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
    const col = db.collection("current")
    return await col.find({}).toArray()
  } 
  
  catch (e) {
    console.error(e)
  } 
  
  finally {
    client.close()
  } 
}
 
const ins = async(data) => {
  let client = new mongo(process.env.PFF, {
    useNewUrlParser: true,
    useUnifiedTopology: true
  })

  try {
    await client.connect()
    let db = client.db("pff")
    let col = db.collection("current")
    
    const insCount = await col.insertMany(data)
    return insCount.insertedCount
  }

  catch (e) {
    console.log(e)
  }
  
  finally {
    client.close()
  }
}

init()
  .then((t) => {
    t.forEach((d) => {
      arr.push(d)
    })
    hist(arr)
      .then((d) => {
        console.log(`${filtered.length} symbols found. Inserting.`)
      })
      .catch((e) => console.log(e))
  })
  .catch((e) => {
    console.log(e)
  })


