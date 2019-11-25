const mongo = require('mongodb').MongoClient

arr = []

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

const hist = async () => {
  let client = new mongo(process.env.PFF, {
    useNewUrlParser: true,
    useUnifiedTopology: true
  })

  try {
    if (arr.length == 0) { 
      return 
    }

    await client.connect()
    const db = client.db("pff")
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

init()
  .then((current) => {
    current.forEach((d) => {
      arr.push(d)
    })
    hist()
      .then((res) => {
        console.log(res)
      })
  })
  .catch((e) => {
    console.log(e)
  })