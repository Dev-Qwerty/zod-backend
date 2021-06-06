const express = require('express')
const http = require('http')
const path = require('path')
const socket = require('socket.io')
const cors = require('cors')

if (process.env.NODE_ENV != 'production') {
    require('dotenv').config()
}

require('./src/config/db')
const router = require('./src/routes/route')

const app = express()

app.use(cors())

app.set('view engine', 'ejs')
app.set('views', './public/views')
app.use(express.static('public'))

router(app)

const server = http.createServer(app)
const io = socket(server)

app.get('/:meetid', (req, res) => {
    res.render('meet')
})

app.set('view engine', 'ejs')
app.use(express.static('public'))

const port = process.env.PORT

app.listen(port, () => {
    console.log("Server running on port " + port + " 🚀")
})