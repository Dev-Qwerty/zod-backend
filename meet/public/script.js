const socket = io()
const peer = new Peer(undefined, {
    host: '/',
    path: '/peerjs',
    port: '8082'
})

peer.on('open', (id) => {
    socket.emit('join-room', ROOM_ID, id)
})

socket.on('user-connected', (userId) => {
    console.log(userId)
})