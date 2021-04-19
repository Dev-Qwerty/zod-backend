const express = require('express')
const { customAlphabet } = require('nanoid')

// Board model
const board = require('../models/boards')

const router = express.Router()
const parseJson = express.json({ extended: true })

// Create 12 digit alphanumeric code for boardId and chardId
const keyStore = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const nanoid = customAlphabet(keyStore, 12)

// @route POST /api/board/new
// @desc Create new board
router
    .route('/new')
    .post( [parseJson], async (req, res) => {
        const { boardName, createdBy, members, projectId, projectName } = req.body

        const boardId = "B" + nanoid()

        // Create three cards(ToDo, Doing, Done) for every board
        const cards = [
            {
                cardId: "C" + nanoid(),
                cardTitle: "To Do"
            },
            {
                cardId: "C" + nanoid(),
                cardTitle: "Doing"
            },
            {
                cardId: "C" + nanoid(),
                cardTitle: "Done"
            }
        ]

        newboard = new board({
            boardId,
            boardName,
            createdBy,
            cards,
            members,
            projectId,
            projectName
        })
        
        await newboard.save()
        res.status(201).send(newboard)
    })

module.exports = router