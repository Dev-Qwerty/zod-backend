const express = require('express')
const { customAlphabet } = require('nanoid')

// Import models
const card = require('../models/card')

const router = express.Router()
const parseJson = express.json({ extended: true })

// Create 12 digit alphanumeric code for cardId
const keyStore = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const nanoid = customAlphabet(keyStore, 16)

// @route POST /api/card/new
// @desc Create new list card
router
    .route('/new')
    .post([parseJson], (req, res) => {
        const { cardName, cardDescription, dueDate, assigned, list, projectId } = req.body

        // Firebase decoded token
        const decodedToken = req.decodedToken

        const cardId = "I" + nanoid()

        createdBy = decodedToken.email

        newCard = new card({
            cardId,
            cardName,
            cardDescription,
            dueDate,
            createdBy,
            assigned,
            list,
            projectId
        })

        newCard.save().then(() => {
            res.status(201).json({
                cardId: newCard.cardId,
                cardName: newCard.cardName,
                cardDescription: newCard.cardDescription,
                dueDate: newCard.dueDate,
                createdBy: newCard.createdBy,
                assigned: newCard.assigned,
                list: newCard.list,
                project: newCard.projectId
            })
        }).catch(error => {
            console.log(error)
            res.status(500).send(error)
        })

    })

router
    .route('/delete')
    .post([parseJson], async (req, res) => {
        const { cardId } = req.body

        const decodedToken = req.decodedToken

        const { email } = decodedToken
        try {
            let doc = await card.findOne({
                cardId, $or: [{ assigned: { $elemMatch: { email } } }, { createdBy: email }]
            })

            if (doc == null) {
                res.status(401).send("Unauthorized user")
            } else {
                let doc = await card.findOneAndDelete({ cardId })
                res.status(200).json({ cardId: doc.cardId, message: "Card deleted" })
            }
        } catch (error) {
            console.log(error)
            res.status(500).send(error)
        }
    })

module.exports = router
