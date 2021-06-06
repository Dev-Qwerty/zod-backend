const socket = io()
const peer = new Peer(undefined, {
    host: '/',
    path: '/peerjs',
    port: '8082'
})

const videoGrid = document.getElementById('video-grid')
const myVideo = document.createElement('video')
myVideo.muted = true
const peers = {}

peer.on('open', (id) => {
    socket.emit('join-room', ROOM_ID, id)
})

navigator.mediaDevices.getUserMedia({
    video: true,
    audio: true
}).then((stream) => {
    addVideoStream(myVideo, stream)

    peer.on('call', (call) => {
        let userId = call.peer
        peers[userId] = call

        call.answer(stream)
        const video = document.createElement('video')
        call.on('stream', (userVideoStream) => {
            addVideoStream(video, userVideoStream)
        })

        call.on('close', () => {
            video.remove()
        })
    })

    socket.on('user-connected', (userId) => {
        setTimeout(connectToNewUser, 1000, userId, stream)
    })

})

socket.on('user-disconnected', (userId) => {
    console.log(userId)
    if (peers[userId]) {
        peers[userId].close()
    }
})

function addVideoStream(video, stream) {
    video.srcObject = stream
    video.addEventListener('loadedmetadata', () => {
        video.play()
    })
    videoGrid.append(video)
}

function connectToNewUser(userId, stream) {
    const call = peer.call(userId, stream)
    const video = document.createElement('video')
    call.on('stream', (userVideoStream) => {
        addVideoStream(video, userVideoStream)
    })

    peers[userId] = call

    call.on('close', () => {
        video.remove()
    })
}