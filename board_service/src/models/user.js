const mongoose = require('mongoose')

const userSchema = new mongoose.Schema({
    fid: String,
    name: String,
    email: String,
    imgUrl: String,
    role: [
        {
            _id: false,
            projectId: String,
            role: String
        }
    ]
})

const user = mongoose.model('users', userSchema)
module.exports = user