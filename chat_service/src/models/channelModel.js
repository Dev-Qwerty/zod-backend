const mongoose = require('mongoose')

const channelSchema = new mongoose.Schema({
    projectid: String,
    channelName: String,
    channelid: String,
    channelAdmin: String,
    members: [String]
})

const channel = mongoose.model('channels', channelSchema);
module.exports = channel