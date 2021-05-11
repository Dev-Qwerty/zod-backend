const express = require('express')
const { customAlphabet } = require('nanoid')

// Import models
const board = require('../models/board')
const user = require('../models/user')

const router = express.Router()
const parseJson = express.json({ extended: true })

// Create 12 digit alphanumeric code for boardId and chardId
const keyStore = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const nanoid = customAlphabet(keyStore, 12)

// @route GET /api/board/:projectid
// @desc Fetch list of all boards
router
    .route('/:projectId')
    .get(async (req, res) => {
        const { projectId } = req.params
        const email = req.decodedToken.email

        try {
            let doc = await user.findOne({
                email, projects: {
                    $elemMatch: { projectId }
                }
            })

            if (doc == null) {
                res.status(401).send("Unauthorized user")
                return
            }

            doc = await board.find({ projectId }, 'boardId boardName type -_id')
            res.status(200).send(doc)

        } catch (error) {
            console.log(error)
            res.status(500).send(error)
        }
    })

// @route POST /api/board/new
// @desc Create new board
router
    .route('/new')
    .post([parseJson], async (req, res) => {
        
        const { boardName, members, type, projectId, projectName } = req.body

        // Firebase decoded token
        const email = req.decodedToken.email

        // Check if the user is inside the project
        let doc = await user.findOne({
            email, projects: {
                $elemMatch: { projectId }
            }
        })

        if (doc == null) {
            res.status(401).send("Unauthorized user")
            return
        }

        const boardId = "B" + nanoid()

        // Add three cards(ToDo, Doing, Done) to every board when created
        const lists = [
            {
                listId: "L" + nanoid(),
                listTitle: "To Do"
            },
            {
                listId: "L" + nanoid(),
                listTitle: "Doing"
            },
            {
                listId: "L" + nanoid(),
                listTitle: "Done"
            }
        ]

        // Push the creater as admin to member array
        members.push({
            email: email,
            isAdmin: true
        })

        newboard = new board({
            boardId,
            boardName,
            lists,
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
        }).catch(error => {
            console.log(error)
            res.status(500).send(error)
        })
    })

// @route POST /api/board/delete
// @desc Delete the board
router
    .route('/delete')
    .post([parseJson], async (req, res) => {

        const { boardId, projectId } = req.body

        // Firebase decoded token
        const email = req.decodedToken.email

        try {
            let doc = await board.findOne({
                boardId, projectId, members: {
                    $elemMatch: { email }
                }
            }, 'members.$ -_id')

            if (doc == null) {
                res.status(401).send("Unauthorized user")
                return
            }

            if (doc.members[0].isAdmin == true) {
                let doc = await board.findOneAndDelete({ boardId, projectId })
                res.status(200).json({ boardId: doc.boardId, message: "Board deleted" })
            } else {
                res.status(401).send("Unathorized user")
            }

        } catch (error) {
            console.log(error)
            res.status(500).send(error)
        }
    })

module.exports = router
