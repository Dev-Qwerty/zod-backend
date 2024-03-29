const express = require('express')

const { createList, updateList, deleteList } = require('./list')
const { createCard, updateCard, deleteCard } = require('./card')
const Card = require('../models/card')
const verifyUser = require('../middlewares/verifyUser')

const router = express.Router()
const parseJson = express.json({ extended: true })

const boardChannel = (namespace, socket, app) => {

    socket.on('joinRoom', (room) => {
        socket.join(room)
    })

    // @route POST /api/:board/list/new
    // @desc Create a new list
    app.post('/api/:board/list/new', [verifyUser, parseJson], async (req, res) => {
        const createdBy = req.decodedToken.email
        const room = req.params.board
        resp = await createList(createdBy, req.body)

        if (resp[0] != "") {
            response = resp[0]
            res.status(201).send('list created')
            namespace.to(room).emit('createList', response)
        } else {
            error = resp[1]
            console.log(error)
            error.message == "Unauthorized user" ? res.status(401).send(error.message) : res.status(500).send(error)
        }
    })

    // @route POST /api/:board/list/update
    // @desc Update the position of board
    app.post('/api/:board/list/update', [verifyUser, parseJson], async (req, res) => {
        const room = req.params.board
        const updatedBy = req.decodedToken.email

        resp = await updateList(updatedBy, req.params.board, req.body)

        if (resp[0] != "") {
            const response = resp[0]
            res.status(201).send('list updated')
            namespace.to(room).emit('updateList', response)
        } else {
            const error = resp[1]
            console.log(error)
            error.message == "Unauthorized user" ? res.status(401).send(error.message) : res.status(500).send(error)
        }
    })

    // @route DELETE /api/:board/list/delete
    // @desc Delete the board
    app.delete('/api/:board/list/delete/:listId', [verifyUser, parseJson], async (req, res) => {
        const room = req.params.board
        const deletedBy = req.decodedToken.email

        resp = await deleteList(deletedBy, req.params)

        if (resp[0] != "") {
            const response = resp[0]
            res.status(201).send('list deleted')
            namespace.to(room).emit('deleteList', response)
        } else {
            const error = resp[1]
            console.log(error)
            error.message == "Unauthorized user" ? res.status(401).send(error.message) : res.status(500).send(error)
        }
    })

    // @route POST /api/:board/card/new
    // @desc Create new list card
    app.post('/api/:board/card/new', [verifyUser, parseJson], async (req, res) => {

        const createdBy = req.decodedToken.email
        const room = req.params.board
        resp = await createCard(createdBy, req.params.board, req.body)

        if (resp[0] != "") {
            const response = resp[0]
            res.status(201).send('card created')
            namespace.to(room).emit('createCard', response)
        } else {
            const error = resp[1]
            console.log(error)
            error.message == "Unauthorized user" ? res.status(401).send(error.message) : res.status(500).send(error)
        }
    })

    // @route POST /api/:board/card/update
    // @desc Update the card
    app.post('/api/:board/card/update', [verifyUser, parseJson], async (req, res) => {
        const email = req.decodedToken.email
        const room = req.params.board

        resp = await updateCard(email, req.params.board, req.body)

        if (resp[0] != "") {
            const response = resp[0]
            res.status(201).send('card updated')
            namespace.to(room).emit('updateCard', response)
        } else {
            const error = resp[1]
            console.log(error)
            error.message == "Unauthorized user" ? res.status(401).send(error.message) : res.status(500).send(error)
        }
    })

    // @route DELETE /api/:board/card/delete/:cardId
    // @desc Delete the card
    app.delete('/api/:board/card/delete/:cardId', [verifyUser, parseJson], async (req, res) => {
        const deletedBy = req.decodedToken.email
        const room = req.params.board

        resp = await deleteCard(deletedBy, req.params)

        if (resp[0] != "") {
            const response = resp[0]
            res.status(201).send('card deleted')
            namespace.to(room).emit('deleteCard', response)
        } else {
            const error = resp[1]
            console.log(error)
            error.message == "Unauthorized user" ? res.status(401).send(error.message) : res.status(500).send(error)
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
