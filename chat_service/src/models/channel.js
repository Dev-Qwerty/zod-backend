const mongoose = require('mongoose')

const memberSchema = new mongoose.Schema({
    name: String,
    fid: String,
    memberid: String,
    imgUrl: String
})

const channelSchema = new mongoose.Schema({
    projectid: String,
    channelName: String,
    channelid: String,
    members: [memberSchema]
})