const mongoose = require('mongoose')

const userSchema = new mongoose.Schema({
    _id: false,
    fid: String,
    name: String,
    email: String,
    imgUrl: String
})

const user = mongoose.model('users', userSchema)
module.exports = user