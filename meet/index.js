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
const Meet = require('./src/models/meet')
const fetchUser = require('./src/utils/fetchUser')

const app = express()

const server = http.createServer(app)
const io = socket(server)
const peerServer = ExpressPeerServer(server)

app.use('/peerjs', peerServer)

app.use(cors())

app.set('view engine', 'ejs')
// app.set('views', './public/views')
app.use(express.static('public'))
app.use(express.json({ extended: true }))

router(app)

app.get('/favicon.ico', (req, res) => res.status(204))

app.post('/leave', async (req, res) => {
    res.status(200).end()
})

app.get('/endmeeting', async (req, res) => {
    const meetId = req.query.meet
    let doc = await Meet.findOne({ meetId }, '-_id meetName')
    if (doc == null) {
        res.end()
    } else {
        res.render('endMeet', { meetName: doc.meetName })
    }
})

app.get('/:meetid', async (req, res) => {
    const meetId = req.params.meetid

    if (req.query.t == undefined) {
        res.render('login')
    } else {
        // Fetch user
        let doc = await fetchUser(req.query.t)

        if (doc == null) {
            res.render('login')
        } else {
            let user = {
                name: doc.name,
                email: doc.email,
                picture: doc.picture
            }

            user = JSON.stringify(user)

            doc = await Meet.findOne({ meetId }, '-_id meetName')
            if (doc == null) {
                res.status(404).end()
            } else {
                process.env.NODE_ENV == 'production' ? meetLink = `https://meet-zode.herokuapp.com/${meetId}` : meetLink = `http://localhost:8082/${meetId}`
                res.render('meet', { meetId, meetLink, meetName: doc.meetName, user })
            }
        }
    }
})

io.on('connection', (socket) => {
    socket.on('join-room', (roomId, userId) => {
        socket.join(roomId)
        socket.to(roomId).emit('user-connected', userId)

        socket.on('disconnect', () => {
            socket.to(roomId).emit('user-disconnected', userId)
        })

        socket.on('new-message', (name, msg) => {
            io.emit('message', name, msg)
        })
    })

})

const port = process.env.PORT

server.listen(port, () => {
    console.log("Server running on port " + port + " 🚀")
})