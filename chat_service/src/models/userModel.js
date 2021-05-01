const mongoose = require('mongoose')

const userSchema = new mongoose.Schema({
    name: String,
    fid: String,
    imgUrl: String,
    email: String,
    projectRole: [
        {
            _id: false,
            projectid: String,
            role: String
        }
    ]
})

const user = mongoose.model('users', userSchema)
module.exports = user