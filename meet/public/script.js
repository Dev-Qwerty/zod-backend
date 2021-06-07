const socket = io()
const peer = new Peer(undefined, {
    secure: true,
    host: 'zode-com-server.herokuapp.com',
    port: 443,
    // host: '/',
    // path: '/peerjs',
    // port: '8082'
})

const videoGrid = document.getElementById('video-grid')
const myVideo = document.createElement('video')
myVideo.muted = true
const peers = {}

peer.on('open', (id) => {
    socket.emit('join-room', MEET_ID, id)
})

navigator.mediaDevices.getUserMedia({
    video: true,
    audio: true
}).then((stream) => {
    myVideoStream = stream
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

const turnOffAudio = () => {
    if (myVideoStream.getAudioTracks()[0].enabled == true) {
        myVideoStream.getAudioTracks()[0].enabled = false
        const showUnmuteBtn = `
            <i class="fas fa-microphone-slash"></i>
        `
        document.querySelector('.main__audio_button').innerHTML = showUnmuteBtn
    } else {
        myVideoStream.getAudioTracks()[0].enabled = true
        const showMuteBtn = `
            <i class="fas fa-microphone"></i>
        `
        document.querySelector('.main__audio_button').innerHTML = showMuteBtn
    }
}

const turnOffVideo = () => {
    if (myVideoStream.getVideoTracks()[0].enabled == true) {
        myVideoStream.getVideoTracks()[0].enabled = false
        const showVideoPauseBtn = `
            <i class="fas fa-video-slash"></i>
         `
        document.querySelector('.main__video_button').innerHTML = showVideoPauseBtn
    } else {
        myVideoStream.getVideoTracks()[0].enabled = true
        const showVideoPlayBtn = `
            <i class="fas fa-video"></i>
         `
        document.querySelector('.main__video_button').innerHTML = showVideoPlayBtn
    }
}

// CHAT FUNCTIONS
const chatForm = document.getElementById('chat-form')

// Message Submit
chatForm.addEventListener('submit', (e) => {
    e.preventDefault()

    // Get message text
    const msg = e.target.elements.msg.value

    if (msg != '') {
        socket.emit('new-message', msg)
    }

})

// Catch msg from server
socket.on('message', (msg) => {
    showMessage(msg)
})

// Show message in DOM
function showMessage(msg) {
    const li = document.createElement('li')
    li.classList.add('message')
    li.innerHTML = `${msg}`
    document.querySelector('.messages').appendChild(li)
}

