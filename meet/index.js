const express = require('express')
const cors = require('cors')

if (process.env.NODE_ENV != 'production') {
    require('dotenv').config()
}

const app = express()

app.set('view engine', 'ejs')
app.use(express.static('public'))

const port = process.env.PORT

app.listen(port, () => {
    console.log("Server running on port " + port + " 🚀")
})