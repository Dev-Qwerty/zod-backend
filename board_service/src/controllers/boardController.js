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
    .post( [parseJson], (req, res) => {
        
        const { boardName, members, type, projectId, projectName } = req.body

        // Firebase decoded token
        const decodedToken = req.decodedToken

        const boardId = "B" + nanoid()

        // Add three cards(ToDo, Doing, Done) to every board when created
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

        // Push the creater as admin to member array
        members.push({
            email: decodedToken.email,
            isAdmin: true
        })

        newboard = new board({
            boardId,
            boardName,
            cards,
            members,
            type,
            projectId,
            projectName
        })

        newboard.save().then( () => {
            res.status(201).json({
                boardId: newboard.boardId,
                boardName: newboard.boardName,
                boardType: newboard.type
            })
        })
    })

module.exports = router
