const mongoose = require('mongoose')

const memeberSchema = new mongoose.Schema({
    _id: false,
    id: String,
    name: String
})

const cardSchema = new mongoose.Schema({
    _id: false,
    cardId: String,
    cardTitle: String,
})

const boardSchema = new mongoose.Schema({
    boardId: String,
    boardName: String,
    createdBy: memeberSchema,
    cards: [cardSchema],
    members: [memeberSchema],
    projectName: String,
    projectId: String
})

const boards = mongoose.model('boards', boardSchema)
module.exports = boards