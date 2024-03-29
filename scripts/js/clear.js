const mongo = require('mongodb').MongoClient

async function init(cli, mdb) {
  try {
    let docs = []
    
    const db = cli.db(mdb)
    const col = db.collection("current")
    
    await col.find({}).toArray().then((arr) => {
      arr.forEach((c) => {
        docs.push(c)
      })
    })

    await db.collection("historical").insertMany(docs)

    return db.collection("current").drop()
  }

  catch (e) {
    console.error(e)
  } 
  
  finally {
    cli.close()
  } 
}

const dbFunc = async(c, dbs) => {
  try {
    return Promise.all(dbs.map(d => init(c, d)))
  }

  catch (e) {
    console.log(e)
  }
}

const client = new mongo(process.env.PFF, {
  useNewUrlParser: true,
  useUnifiedTopology: true
}).connect()
  .then((c) => {
    dbFunc(c, ["cef", "pff"])
    .then(res => {
      console.log(res)
    })
  })

