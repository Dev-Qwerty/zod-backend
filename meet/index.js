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

app.get('/getfbc', (req, res) => {
    let fbc = {
        apiKey: process.env.APIKEY,
        authDomain: process.env.AUTHDOMAIN,
        projectId: process.env.PROJECTID,
        storageBucket: process.env.STORAGEBUCKET,
        messagingSenderId: process.env.MESSAGINGSENDERID,
        appId: process.env.APPID
    }

    const fbcStr = JSON.stringify(fbc)
    fbc = Buffer.from(fbcStr).toString("base64")

    res.status(200).json({ fbc: fbc })
})

app.get('/postmeet', async (req, res) => {
    const meetId = req.query.meet
    const end = req.query.end
    let doc = await Meet.findOne({ meetId }, '-_id meetName')
    if (doc == null) {
        res.end()
    } else {
        if (end != undefined) {
            res.render('endMeet', { meetName: doc.meetName })
        } else {
            res.render('leaveMeet', { meetName: doc.meetName })
        }
    }
})

app.get('/:meetid', async (req, res) => {
    const meetId = req.params.meetid

    const meet = await Meet.findOne({ meetId })
    if (meet != null) {
        if (req.query.t == undefined) {
            res.render('login', { meetId, meetName: meet.meetName })
        } else {
            // Fetch user
            let doc = await fetchUser(req.query.t)

            if (doc == null) {
                res.render('login', { meetId, meetName: meet.meetName })
            } else {
                const createdBy = doc.email
                const isCreator = await Meet.findOne({ createdBy })
                let user = {
                    name: doc.name,
                    email: doc.email,
                    picture: doc.picture,
                    host: isCreator != null ? true : false
                }

                user = JSON.stringify(user)

                process.env.NODE_ENV == 'production' ? meetLink = `https://meet-zode.herokuapp.com/${meetId}` : meetLink = `http://localhost:8082/${meetId}`
                res.render('meet', { meetId, meetLink, meetName: meet.meetName, user })
            }
        }
    } else {
        res.render('noMeet')
    }
})

io.on('connection', (socket) => {
    socket.on('join-room', (roomId, userId) => {
        socket.join(roomId)
        socket.broadcast.to(roomId).emit('user-connected', userId)

        socket.on('disconnect', () => {
            socket.broadcast.to(roomId).emit('user-disconnected', userId)
        })

        socket.on('new-message', (name, msg) => {
            io.to(roomId).emit('message', name, msg)
        })

        socket.on('end-meeting', () => {
            io.to(roomId).emit('end')
        })
    })

})

const port = process.env.PORT

server.listen(port, () => {
    console.log("Server running on port " + port + " 🚀")
})