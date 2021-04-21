const mongoose = require('mongoose')

const boardSchema = new mongoose.Schema({
    boardId: String,
    boardName: String,
    cards: [
        {
            _id: false,
            cardId: String,
            cardTitle: String
        }
    ],
    members: [
        {
            _id: false,
            email: String,
            isAdmin: {
                type: Boolean,
                default: false
            }
        }
    ],
    projectId: String,
    projectName: String
})

const boards = mongoose.model('boards', boardSchema)
module.exports = boards