const mongoose = require('mongoose')

const meetSchema = new mongoose.Schema({
    meetId: String,
    meetUrl: String,
    createdBy: String,
    projectId: String,
    members: [
        {
            _id: false,
            email: String,
        }
    ]
})

const meets = mongoose.model('meets', meetSchema)
module.exports = meets