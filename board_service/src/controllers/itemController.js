const express = require('express')
const { customAlphabet } = require('nanoid')

// Import models
const item = require('../models/items')

const router = express.Router()
const parseJson = express.json({ extended: true })

// Create 12 digit alphanumeric code for itemId
const keyStore = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const nanoid = customAlphabet(keyStore, 16)

// @route POST /api/item/new
// @desc Create new card item
router
    .route('/new')
    .post([parseJson], (req, res) => {
        const { itemName, itemDescription, dueDate, assigned, card, projectId } = req.body

        // Firebase decoded token
        const decodedToken = req.decodedToken

        const itemId = "I" + nanoid()

        createdBy = decodedToken.email

        newItem = new item({
            itemId,
            itemName,
            itemDescription,
            dueDate,
            createdBy,
            assigned,
            card,
            projectId
        })

        newItem.save().then(() => {
            res.status(201).json({
                itemId: newItem.itemId,
                itemName: newItem.itemName,
                itemDescription: newItem.itemDescription,
                dueDate: newItem.dueDate,
                createdBy: newItem.createdBy,
                assigned: newItem.assigned,
                card: newItem.card,
                project: newItem.projectId
            })
        }).catch(error => {
            console.log(error)
            res.status(500).send(error)
        })

    })

module.exports = router
