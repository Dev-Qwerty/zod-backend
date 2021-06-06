const express = require('express')
const cors = require('cors')

if (process.env.NODE_ENV != 'production') {
    require('dotenv').config()
}

require('./src/config/db')
const router = require('./src/routes/route')

const app = express()

app.use(cors())

router(app)

app.set('view engine', 'ejs')
app.use(express.static('public'))

const port = process.env.PORT

app.listen(port, () => {
    console.log("Server running on port " + port + " 🚀")
})