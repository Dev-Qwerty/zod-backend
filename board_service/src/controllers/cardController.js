const express = require('express')
const { customAlphabet } = require('nanoid')

// Import models
const Card = require('../models/card')
const List = require('../models/list')

const router = express.Router()
const parseJson = express.json({ extended: true })

// Create 12 digit alphanumeric code for cardId
const keyStore = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const nanoid = customAlphabet(keyStore, 16)

const boardChannel = (namespace, socket, app) => {

    socket.on('joinRoom', (room) => {
        socket.join(room)
    })

    // @route PSOT /api/list/new
    // @desc Create a new list
    app.post('/api/:board/list/new', [parseJson], async (req, res) => {
        try {
            const room = req.params.board
            const { title, pos, boardId } = req.body

            const email = socket.email

            const newlist = new List({
                listId: "L" + nanoid(),
                title,
                pos,
                createdBy: email,
                boardId
            })

            await newlist.save()
            const response = {
                listId: newlist.listId,
                title: newlist.title,
                pos: newlist.pos
            }
            res.status(201).send('list created')
            namespace.to(room).emit('createList', response)
        } catch (error) {
            console.log(error)
            res.status(500).send(error)
        }
    })

    // @route POST /api/card/new
    // @desc Create new list card
    app.post('/api/:board/card/new', [parseJson], async (req, res) => {
        try {
            const room = req.params.board
            const { cardName, cardDescription, dueDate, pos, assigned, list, projectId } = req.body

            const cardId = "I" + nanoid()

            createdBy = socket.email

            newCard = new Card({
                cardId,
                cardName,
                cardDescription,
                dueDate,
                pos,
                createdBy,
                assigned,
                list,
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

            res.status(201).send('card created')
            namespace.to(room).emit('createCard', response)

        } catch (error) {
            console.log(error)
            res.status(500).send(error)
        }
    })


    router
        .route('/delete')
        .post([parseJson], async (req, res) => {
            const { cardId } = req.body

            const decodedToken = req.decodedToken

            const { email } = decodedToken
            try {
                let doc = await Card.findOne({
                    cardId, $or: [{ assigned: { $elemMatch: { email } } }, { createdBy: email }]
                })

                if (doc == null) {
                    res.status(401).send("Unauthorized user")
                } else {
                    let doc = await Card.findOneAndDelete({ cardId })
                    res.status(200).json({ cardId: doc.cardId, message: "Card deleted" })
                }
            } catch (error) {
                console.log(error)
                res.status(500).send(error)
            }
        })
}


// module.exports = router
module.exports = boardChannel
