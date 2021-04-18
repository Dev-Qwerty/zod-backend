const mongoose = require('mongoose')

const memeberSchema = new mongoose.Schema({
    name: String,
    memberId: String
})

const cardSchema = new mongoose.Schema({
    cardId: String,
    cardTitle: String,
})

const boardSchema = new mongoose.Schema({
    boardName: String,
    cards: [cardSchema],
    members: [memeberSchema],
    projectName: String,
    projectId: String
})