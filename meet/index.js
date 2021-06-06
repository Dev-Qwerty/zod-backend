const express = require('express')
const http = require('http')
const path = require('path')
const socket = require('socket.io')
// const { ExpressPeerServer } = require('peer')
const cors = require('cors')

if (process.env.NODE_ENV != 'production') {
    require('dotenv').config()
}

require('./src/config/db')
const router = require('./src/routes/route')
const Meet = require('./src/models/meet')

const app = express()

const server = http.createServer(app)
const io = socket(server)
// const peerServer = ExpressPeerServer(server)

// app.use('/peerjs', peerServer)

app.use(cors())

app.set('view engine', 'ejs')
// app.set('views', './public/views')
app.use(express.static('public'))

router(app)

app.get('/favicon.ico', (req, res) => res.status(204))

app.get('/:meetid', async (req, res) => {
    const meetId = req.params.meetid
    let doc = await Meet.findOne({ meetId }, '-_id meetName')
    if (doc == null) {
        res.status(404).end()
    } else {
        process.env.NODE_ENV == 'production' ? meetLink = `https://meet-zode.herokuapp.com/${meetId}` : meetLink = `http://localhost:8082/${meetId}`
        res.render('meet', { meetId, meetLink, meetName: doc.meetName })
    }
})

io.on('connection', (socket) => {
    socket.on('join-room', (roomId, userId) => {
        socket.join(roomId)
        socket.to(roomId).emit('user-connected', userId)

        socket.on('disconnect', () => {
            socket.to(roomId).emit('user-disconnected', userId)
        })
    })

})

const port = process.env.PORT

server.listen(port, () => {
    console.log("Server running on port " + port + " 🚀")
})