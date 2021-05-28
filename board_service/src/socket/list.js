const { customAlphabet } = require('nanoid')

// Import models
const List = require('../models/list')

// Create 12 digit alphanumeric code for cardId
const keyStore = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const nanoid = customAlphabet(keyStore, 16)

// @route PSOT /api/:board/list/new
// @desc Create a new list
async function createList(createdBy, data) {

    try {
        const { title, pos, boardId } = data

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

// module.exports = router
module.exports = createList
