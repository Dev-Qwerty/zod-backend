const express = require('express')
const http = require('http')
const path = require('path')
const socket = require('socket.io')
const { ExpressPeerServer } = require('peer')
const cors = require('cors')

if (process.env.NODE_ENV != 'production') {
    require('dotenv').config()
}

require('./src/config/db')
const router = require('./src/routes/route')

const app = express()

const server = http.createServer(app)
const io = socket(server)
const peerServer = ExpressPeerServer(server)

app.use('/peerjs', peerServer)

app.use(cors())

app.set('view engine', 'ejs')
// app.set('views', './public/views')
app.use(express.static('public'))

router(app)

app.get('/:meetid', (req, res) => {
    res.render('meet')
})

io.on('connection', (socket) => {
    socket.on('join-room', (roomId, userId) => {
        socket.join(roomId)
        socket.to(roomId).emit('user-connected', userId)
    })
})

const port = process.env.PORT

server.listen(port, () => {
    console.log("Server running on port " + port + " 🚀")
})