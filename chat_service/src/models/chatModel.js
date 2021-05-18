const mongoose = require('mongoose')

const chatSchema = new mongoose.Schema({
    projectid: String,
    channelid: String,
    sender: {
        name: String,
        email: String
    },
    date: {
        type: Date,
        default: Date.now()
    },
    type: String,
    message: String
})

const chat = mongoose.model('chats', chatSchema);
module.exports = chat