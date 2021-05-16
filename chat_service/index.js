const express = require('express')
const http = require('http')
const socketio = require('socket.io')
require('dotenv').config()

require('./src/config/db')
const firebaseAuth = require('./src/config/firebase')

const VerifyUser = require('./src/middlewares/verifyuser')

// const kafkaConsumer = require('./src/messageQueue/consumer')
// kafkaConsumer()

const app = express()
const server = http.createServer(app)

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
    require('./src/socket/channelMessage')(projectSpace, socket)
})

server.listen(process.env.PORT, () => {
    console.log(`Server running on port ${process.env.PORT}`)
})