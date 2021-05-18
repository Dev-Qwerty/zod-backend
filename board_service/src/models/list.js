const mongoose = require('mongoose')

const listSchema = new mongoose.Schema({
    listId: String,
    title: String,
    index: Number,
    createdBy: String,
    boardId: String
})

const lists = mongoose.model('lists', listSchema)
module.exports = lists
