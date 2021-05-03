const mongoose = require('mongoose')

const itemSchema = new mongoose.Schema({
    itemId: String,
    itemName: String,
    itemDescription: String,
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
    card: {
        _id: false,
        cardId: String,
        cardTitle: String
    },
    projectId: String
})

const items = mongoose.model('items', itemSchema)
module.exports = items
