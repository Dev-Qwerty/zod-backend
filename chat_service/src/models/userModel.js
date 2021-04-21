const mongoose = require('mongoose')

const userSchema = new mongoose.Schema({
    name: String,
    fid: String,
    imgUrl: String,
    email: String
})

const user = mongoose.model('users', userSchema)
module.exports = user