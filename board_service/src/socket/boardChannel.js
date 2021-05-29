const express = require('express')

const { createList, updateList } = require('./list')
const createCard = require('./card')
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
            res.status(500).send(error)
        }
    })

    // @route PUT /api/:board/list/pos
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
            res.status(500).send(error)
        }
    })

    // @route POST /api/:board/card/new
    // @desc Create new list card
    app.post('/api/:board/card/new', [parseJson], async (req, res) => {

        const createdBy = req.decodedToken.email
        const room = req.params.board
        resp = await createCard(createdBy, req.body)

        if (resp[0] != "") {
            const response = resp[0]
            res.status(201).send('card created')
            namespace.to(room).emit('createCard', response)
        } else {
            const error = resp[1]
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
