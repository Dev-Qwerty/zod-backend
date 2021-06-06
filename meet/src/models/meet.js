const mongoose = require('mongoose')

const meetSchema = new mongoose.Schema({
    meetId: String,
    meetName: String,
    meetUrl: String,
    createdBy: String,
    projectId: String,
    members: Array
})

const meets = mongoose.model('meets', meetSchema)
module.exports = meets