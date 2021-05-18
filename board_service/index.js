const express = require('express')
const socket = require('socket.io')
const http = require('http')
require('dotenv').config()

require('./src/config/db')
const router = require('./src/routes/route')
// require('./src/messageQueues/consumer').kafkaConsumer
const auth = require('./src/config/firebase')

const app = express()

router(app)

const server = http.createServer(app)
const io = socket(server, {
    cors: {
        origin: "http://localhost:8080",
        methods: ['GET', 'POST', 'PUT', 'DELETE']
    }
})

const boardNamespace = io.of(/.*$/)

boardNamespace.use((socket, next) => {
    auth.verifyIdToken(socket.handshake.auth.token)
        .then((decodedToken) => {
            socket.email = decodedToken.email
            next()
        }).catch((error) => {
            next(new Error('Unauthorized user'))
        })
})

boardNamespace.on('connection', (socket) => {
    const namespace = socket.nsp
    require('./src/controllers/cardController')(namespace, socket, app)
})

server.listen(3000, () => {
    console.log('server running...')
});
