const mongoose = require('mongoose')

const channelSchema = new mongoose.Schema({
    projectid: String,
    channelName: String,
    channelid: String,
    channelAdmin: String,
    description: String,
    members: [
        {
            _id: false,
            email: String,
            isAdmin: {
                type: Boolean,
                default: false
            }
        }
    ]
})

const channel = mongoose.model('channels', channelSchema);
module.exports = channel