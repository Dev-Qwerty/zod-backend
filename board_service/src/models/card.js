const mongoose = require('mongoose')

const cardSchema = new mongoose.Schema({
    cardId: String,
    cardName: String,
    cardDescription: String,
    dueDate: String,
    pos: Number,
    assigned: [
        {
            _id: false,
            name: String,
            email: String,
            imgUrl: String
        }
    ],
    createdBy: String,
    listId: String,
    projectId: String
})

const cards = mongoose.model('cards', cardSchema)
module.exports = cards
