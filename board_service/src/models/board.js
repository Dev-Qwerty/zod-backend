const mongoose = require('mongoose')

const boardSchema = new mongoose.Schema({
    boardId: String,
    boardName: String,
    lists: [
        {
            _id: false,
            listId: String,
            listTitle: String
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
    type: String,
    projectId: String,
    projectName: String
})

const boards = mongoose.model('boards', boardSchema)
module.exports = boards
