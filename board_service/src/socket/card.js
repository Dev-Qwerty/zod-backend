const { customAlphabet } = require('nanoid')

// Import models
const Card = require('../models/card')
const Board = require('../models/board')

// Create 12 digit alphanumeric code for cardId
const keyStore = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const nanoid = customAlphabet(keyStore, 16)

// @route POST /api/:board/card/new
// @desc Create new list card
async function createCard(createdBy, boardId, data) {

    try {
        const { cardName, cardDescription, dueDate, pos, assigned, listId } = data

        const email = createdBy
        // Check if the user creating card is inside the board
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
            boardId
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
            listId: newCard.listId,
        }

        return [response, ""]

    } catch (error) {
        console.log(error)
        return ["", error]
    }
}

async function updateCard(updatedBy, boardId, data) {
    try {
        const { cardId, cardName, cardDescription, pos, dueDate, listId } = data
        const email = updatedBy
        // Check if the user updating card is inside the board
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

        let update = {}

        cardName != undefined ? update.cardName = cardName : ""
        cardDescription != undefined ? update.cardDescription = cardDescription : ""
        pos != undefined ? update.pos = pos : ""
        dueDate != undefined ? update.dueDate = dueDate : ""
        listId != undefined ? update.listId = listId : ""

        const card = await Card.findOneAndUpdate({ cardId }, update, { new: true })

        const response = {
            cardId: card.cardId,
            cardName: card.cardName,
            cardDescription: card.cardDescription,
            dueDate: card.dueDate,
            pos: card.pos,
            createdBy: card.createdBy,
            assigned: card.assigned,
            listId: card.listId,
            boardId: card.boardId
        }

        return [response, ""]

    } catch (error) {
        console.log(error)
        return ["", error]
    }
}

async function deleteCard(deletedBy, data) {
    try {
        console.log(data)
        const { board, cardId } = data

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

        doc = await Card.findOneAndDelete({ cardId })

        const response = {
            cardId: doc.cardId
        }

        return [response, ""]

    } catch (error) {
        console.log(error)
        return ["", error]
    }
}

// module.exports = router
module.exports = { createCard, updateCard, deleteCard }
