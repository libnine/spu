const axios = require('axios')
const mongo = require('mongodb').MongoClient

let arr = []
let arr_cef = []

async function init(d) {
  let client = new mongo(process.env.PFF, {
    useNewUrlParser: true,
    useUnifiedTopology: true
  })

  try {
    await client.connect()
    const db = client.db(d)
    const col = db.collection("symbols")
    return await col.find({}).sort({"symbol": 1}).toArray()
  }

  catch (e) {
    console.error(e)
  }

  finally {
    client.close()
  }
}

// 

const ins = async(d, data) => {
  let client = new mongo(process.env.PFF, {
    useNewUrlParser: true,
    useUnifiedTopology: true
  })

  try {
    await client.connect()
    console.log(`${d} insert starting.`)

    let db = client.db(d)
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

const scraped = async (ticker, type) => {
  try {
    let res = await axios.get(`https://www.marketwatch.com/investing/${type}/${ticker}`)

    dump = JSON.parse(/\s+\<script\stype\=\"application\/[a-zA-Z+]+\"\>\s+?({[\s\S]*})\s+?\<\/script\>/g.exec(res.data)[1])
    re_52wk_low = /52\s\w+\s\w+\<\/[\w]+\>[\s\<\w\=\"\w+]+\>(.+?)\s\-\s[\d\.]+\<\//g
    re_52wk_high = /52\s\w+\s\w+\<\/[\w]+\>[\s\<\w\=\"\w+]+\>[\d\.]+\s\-\s(.+?)\<\//g
    re_chg = /change\-\-point\-\-q\"\>(.+?)\</g
    re_chg_pct = /change\-\-percent\-\-q\"\>(.+?)\</g
    re_day_high = /Day\sRange\<\/[a-z]+\>\s+\<span\s[a-z]+\=\"[a-z\s\_]+\"\>[0-9\.]+\s\-\s(.+?)\</g
    re_day_low = /Day\sRange\<\/[a-z]+\>\s+\<span\s[a-z]+\=\"[a-z\s\_]+\"\>(.+?)\s\-\s[0-9\.]+\</g
    re_div = /Dividend\<\/[a-z]+\>\s+\<span\s[a-z]+\=\"[a-z\s\_]+\"\>\$(.+?)\</g
    re_exdiv = /Ex\-Dividend\sDate\<\/[a-z]+\>\s+\<span\s[a-z]+\=\"[a-z\s\_]+\"\>(.+?)\</g
    re_price = /sup\>\s+\<span\sclass="value"\>(.+?)\<\/span>/g
    re_vol = /volume\slast\-value\"\>\s+(.+?)K?\s+/g
    re_vol_rel = /vs\-average\"\>(.+?)\%\</g
    re_yield = /Yield\<\/[a-z]+\>\s+\<span\s[a-z]+\=\"[a-z\s\_]+\"\>(.+?)%\</g

    return ({
      "52wk_low": parseFloat(re_52wk_low.exec(res.data)[1]),
      "52wk_high": parseFloat(re_52wk_high.exec(res.data)[1]),
      "chg": parseFloat(dump.priceChange),
      "chg_pct": parseFloat(dump.priceChangePercent.replace("%", "")),
      "date": new Date(Date.now()).toISOString(),
      "day_high": parseFloat(re_day_high.exec(res.data)[1]),
      "day_low": parseFloat(re_day_low.exec(res.data)[1]),
      "div": parseFloat(re_div.exec(res.data)[1]),
      "ex_div": re_exdiv.exec(res.data)[1],
      "last": parseFloat(dump.price),
      "name": dump.name,
      "quoteTime": dump.quoteTime,
      "relative_volume": (parseFloat(re_vol_rel.exec(res.data)[1]) / 100),
      "ticker": ticker,
      "volume": parseFloat(re_vol.exec(res.data)[1]),
      "yield": parseFloat(re_yield.exec(res.data)[1])
    })
  }
  catch (e) {console.log(e)}
}

async function go(dump, type) {
  try {
    return Promise.all(dump.map(a => scraped(a, type)))
  }
  catch (e) {
    return e
  }
}

init("cef")
  .then((t) => {
    t.forEach((d) => {
      arr_cef.push(d.symbol.toLowerCase())
    })
    go(arr_cef, "fund")
      .then((data) => {
        let filtered = data.filter((el) => {
          return el != null
        })
        console.log(`${filtered.length} symbols found. Inserting.`)
        ins("cef", filtered)
          .then((res) => {
            console.log(res)
          })
      })
      .catch((e) => console.log(e))
  })
  .catch((e) => {
    console.log(e)
  })

init("pff")
  .then((t) => {
    t.forEach((d) => {
      arr.push(d.symbol.toLowerCase().replace(".p", ".pr"))
    })
    go(arr, "stock")
      .then((d) => {
        let filtered = d.filter((el) => {
          return el != null
        })
        console.log(`${filtered.length} symbols found. Inserting.`)
        ins("pff", filtered)
          .then((res) => {
            console.log(res)
          })
      })
      .catch((e) => console.log(e))
  })
  .catch((e) => {
    console.log(e)
  })