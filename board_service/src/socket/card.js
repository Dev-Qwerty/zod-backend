const { customAlphabet } = require('nanoid')

// Import models
const Card = require('../models/card')

// Create 12 digit alphanumeric code for cardId
const keyStore = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const nanoid = customAlphabet(keyStore, 16)

// @route POST /api/:board/card/new
// @desc Create new list card
async function createCard(createdBy, data) {

    try {
        const { cardName, cardDescription, dueDate, pos, assigned, listId, projectId } = data

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

        const cardId = "I" + nanoid()

        newCard = new Card({
            cardId,
            cardName,
            cardDescription,
            dueDate,
            pos,
            createdBy,
            assigned,
            listId,
            projectId
        })

        await newCard.save()
        const response = {
            cardId: newCard.cardId,
            cardName: newCard.cardName,
            cardDescription: newCard.cardDescription,
            dueDate: newCard.dueDate,
            pos: newCard.pos,
            createdBy: newCard.createdBy,
            assigned: newCard.assigned,
            list: newCard.list,
            project: newCard.projectId
        }

        return [response, ""]

    } catch (error) {
        console.log(error)
        return ["", error]
    }
}

// module.exports = router
module.exports = createCard
