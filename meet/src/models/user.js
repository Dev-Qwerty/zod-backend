const mongoose = require('mongoose')

const userSchema = new mongoose.Schema({
    name: String,
    fid: String,
    email: String,
    imgUrl: String,
    projects: [
        {
            _id: false,
            projectId: String,
            role: String
        }
    ]
})

const user = mongoose.model('users', userSchema)
module.exports = user