const { customAlphabet } = require('nanoid')

// Import models
const List = require('../models/list')
const Board = require('../models/board')
const Card = require('../models/card')

// Create 12 digit alphanumeric code for cardId
const keyStore = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const nanoid = customAlphabet(keyStore, 16)

// @route PSOT /api/:board/list/new
// @desc Create a new list
async function createList(createdBy, data) {

    try {
        const { title, pos, boardId } = data

        const email = createdBy
        // Check if the user updating list is inside the board
        let doc = await Board.findOne({
            boardId, members: {
                $elemMatch: { email }
            }
        })

        if (doc == null) {
            const error = {
                message: "Unauthorized user"
            }
            return ["", error]
        }

        const newlist = new List({
            listId: "L" + nanoid(),
            title,
            pos,
            createdBy,
            boardId
        })

        await newlist.save()
        const response = {
            listId: newlist.listId,
            title: newlist.title,
            pos: newlist.pos
        }
        return [response, ""]
    } catch (error) {
        console.log(error)
        return ["", error]
    }
}

async function updateList(updatedBy, boardId, data) {
    try {
        const { listId, title, pos } = data
        const email = updatedBy
        // Check if the user updating list is inside the board
        let doc = await Board.findOne({
            boardId, members: {
                $elemMatch: { email }
            }
        })

        if (doc == null) {
            const error = {
                message: "Unauthorized user"
            }
            return ["", error]
        }

        if (title != undefined) {
            doc = await List.findOneAndUpdate({ listId }, { title }, { new: true })
        }

        if (pos != undefined) {
            doc = await List.findOneAndUpdate({ listId }, { pos }, { new: true })
        }

        const response = {
            title: doc.title,
            listId: doc.listId,
            pos: doc.pos,
            boardId: doc.boardId
        }
        return [response, ""]

    } catch (error) {
        console.log(error)
        return ["", error]
    }
}

async function deleteList(deletedBy, data) {
    try {
        const { board, listId } = data

        const email = deletedBy
        // Check if the user updating list is inside the board
        let doc = await Board.findOne({
            boardId: board, members: {
                $elemMatch: { email }
            }
        })

        if (doc == null) {
            const error = {
                message: "Unauthorized user"
            }
            return ["", error]
        }

        doc = await List.findOneAndDelete({ listId })
        await Card.deleteMany({ listId })

        const response = {
            listId: doc.listId
        }
        return [response, ""]
    } catch (error) {
        console.log(error)
        return ["", error]
    }
}

// module.exports = router
module.exports = { createList, updateList, deleteList }
