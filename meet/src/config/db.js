const mongoose = require('mongoose')

mongoose.connect(process.env.MONGOURI, { useNewUrlParser: true, useUnifiedTopology: true, useFindAndModify: true })
    .catch((error) => {
        console.error(`Failed to connect to db: ${error}`)
    })

const db = mongoose.connection

db.on('open', () => {
    console.log("Connected to DB.")
})

// Handle error after connection
db.on('error', (error) => {
    console.error(`DB connection failed: ${error}`)
})