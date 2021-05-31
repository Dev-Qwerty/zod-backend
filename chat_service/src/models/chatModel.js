const mongoose = require('mongoose')

const chatSchema = new mongoose.Schema({
    messageid: String,
    projectid: String,
    channelid: String,
    author: {
        imgUrl: String,
        name: String,
        email: String
    },
    ts: Number,
    content: String
})

const chat = mongoose.model('chats', chatSchema);
module.exports = chat