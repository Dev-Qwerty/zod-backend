const express = require('express')
const http = require('http')
const socketio = require('socket.io')
const { customAlphabet } = require('nanoid');
const cors = require('cors')

if (process.env.NODE_ENV != "production") {
    require('dotenv').config()
}

const channelModel = require('./src/models/channelModel')
const userModel = require('./src/models/userModel')

const alphabet = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
const nanoid = customAlphabet(alphabet, 12);

require('./src/config/db')
const firebaseAuth = require('./src/config/firebase')

const VerifyUser = require('./src/middlewares/verifyuser')

// const kafkaConsumer = require('./src/messageQueue/consumer')
// kafkaConsumer()

const app = express()
const server = http.createServer(app)

app.use(cors())


app.post('/api/channel/everyone/new', [parseJson], async (req, res) => {
    console.log(req.body)
    const member = {}
    const newChannel = new channelModel({
        projectid: req.body.projectid,
        channelName: 'everyone',
        channelid: nanoid(),
    })
    member.email = req.body.member.email
    member.isAdmin = false
    newChannel.members.push(member)
    newChannel.save()
        .then(doc => {
            console.log(doc)
            res.end()
        })
        .catch(error => {
            console.log(error)
            res.end()
        })
    const { name, fid, imgUrl, email } = req.body.member
    projectRole = []
    projectRole.push({
        projectid: req.body.projectid,
        role: req.body.member.role
    })
    userModel.findOneAndUpdate({ email }, { name, fid, email, imgUrl, $push: { projectRole } }, { upsert: true, new: true })
        .then(doc => {
            console.log(doc)
            res.end()
        })
        .catch(error => {
            console.log(error)
            res.end()
        })
})

app.post('/api/channel/acceptinvite', [parseJson], async (req, res) => {
    console.log(req.body)
    const { name, fid, imgUrl, email } = req.body.member
    await channelModel.findOneAndUpdate({ projectid: req.body.projectid, channelName: 'everyone' }, { $push: { members: { isAdmin: false, email } } })
    projectRole = []
    projectRole.push({
        projectid: req.body.projectid,
        role: req.body.member.role
    })
    await userModel.findOneAndUpdate({ email }, { name, fid, email, imgUrl, $push: { projectRole } }, { upsert: true, new: true })
    res.end()
})

app.use(VerifyUser)

//Initialize Routes
require('./src/routes/route')(app)

// Initialize socket
const io = socketio(server, {
    cors: {
        origin: ["http://localhost:8080", "http://localhost:3000"],
    }
})

const projectSpaces = io.of(/.*$/)
projectSpaces.use((socket, next) => {
    const idToken = socket.handshake.auth.token
    firebaseAuth.verifyIdToken(idToken)
        .then((decodedToken) => {
            socket.email = decodedToken.email
            next()
        })
        .catch((error) => {
            console.log('VerifyUser: ', error)
            next(new Error('Authentication error'));
        })
})

projectSpaces.on("connection", socket => {
    const projectSpace = socket.nsp
    socket.emit("connection", "connected to server")
    require('./src/socket/channelMessage')(projectSpace, socket, app)
})

server.listen(process.env.PORT, () => {
    console.log(`Server running on port ${process.env.PORT}`)
})