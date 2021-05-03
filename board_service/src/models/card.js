const mongoose = require('mongoose')

const cardSchema = new mongoose.Schema({
    cardId: String,
    cardName: String,
    cardDescription: String,
    dueDate: String,
    assigned: [
        {
            _id: false,
            name: String,
            email: String,
            imgUrl: String
        }
    ],
    createdBy: String,
    list: {
        _id: false,
        listId: String,
        listTitle: String
    },
    projectId: String
})

const cards = mongoose.model('cards', cardSchema)
module.exports = cards
