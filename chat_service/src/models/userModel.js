const mongoose = require('mongoose')

const userSchema = new mongoose.Schema({
    _id: false,
    name: String,
    fid: String,
    memberid: String,
    imgUrl: String,
    role: String,
    email: String
})

const user = mongoose.model('users', userSchema)
module.exports = user