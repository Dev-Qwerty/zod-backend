const mongoose = require('mongoose')

const boardSchema = new mongoose.Schema({
    boardId: String,
    boardName: String,
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
    createdBy: String,
    projectId: String,
    projectName: String
})

const boards = mongoose.model('boards', boardSchema)
module.exports = boards
