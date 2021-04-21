const mongoose = require('mongoose')

const memberSchema = new mongoose.Schema({
    _id: false,
    name: String,
    fid: String,
    memberid: String,
    imgUrl: String,
    role: String,
    email: String,
    channelAdmin: Boolean
})

const channelSchema = new mongoose.Schema({
    projectid: String,
    channelName: String,
    channelid: String,
    members: [memberSchema]
})

const channel = mongoose.model('channels', channelSchema);
module.exports = channel