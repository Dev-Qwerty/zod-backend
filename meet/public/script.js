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

navigator.mediaDevices.getUserMedia({
    video: true,
    audio: true
}).then((stream) => {
    addVideoStream(myVideo, stream)
})

peer.on('open', (id) => {
    socket.emit('join-room', ROOM_ID, id)
})

socket.on('user-connected', (userId) => {
    console.log(userId)
})

function addVideoStream(video, stream) {
    video.srcObject = stream
    video.addEventListener('loadedmetadata', () => {
        video.play()
    })
    videoGrid.append(video)
}