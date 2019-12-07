const mongo = require('mongodb').MongoClient

arr = []

async function init(mdb) {
  let client = new mongo(process.env.PFF, {
    useNewUrlParser: true,
    useUnifiedTopology: true
  })

  try {
    arr.splice(0, arr.length)
    await client.connect()
    const db = client.db(mdb)
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

const hist = async (mdb) => {
  let client = new mongo(process.env.PFF, {
    useNewUrlParser: true,
    useUnifiedTopology: true
  })

  try {
    if (arr.length == 0) { 
      return 
    }

    await client.connect()
    const db = client.db(mdb)
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
  .then((current) => {
    current.forEach((d) => {
      arr.push(d)
  })
  hist("pff")
    .then((res) => {
      console.log(res)
  })
init("cef")
  .then((current) => {
    current.forEach((d) => {
      arr.push(d)
  })
  hist("cef")
    .then((res) => {
      console.log(res)
    })
  })
})
