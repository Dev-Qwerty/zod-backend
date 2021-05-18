const express = require('express')
const { customAlphabet } = require('nanoid')

// Import models
const Board = require('../models/board')
const User = require('../models/user')
const List = require('../models/list')

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
            let doc = await User.findOne({
                email, projects: {
                    $elemMatch: { projectId }
                }
            })

            if (doc == null) {
                res.status(401).send("Unauthorized user")
                return
            }

            doc = await Board.find({ projectId }, 'boardId boardName type -_id')
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
        try {
            const { boardName, type, projectId, projectName } = req.body
            let members = req.body.members

            // Firebase decoded token
            const email = req.decodedToken.email

            // Check if the user is inside the project
            let doc = await User.findOne({
                email, projects: {
                    $elemMatch: { projectId }
                }
            })

            if (doc == null) {
                res.status(401).send("Unauthorized user")
                return
            }

            const boardId = "B" + nanoid()

            if (type == "public") {
                // Push the creator as admin to member array(public boards)
                members.push({
                    email: email,
                    isAdmin: true
                })
            } else {
                // Set creator as admin to member array(private boards)
                members = {
                    email: email,
                    isAdmin: true
                }
            }

            let newboard = new Board({
                boardId,
                boardName,
                members,
                type,
                projectId,
                projectName
            })

            await newboard.save()

            // Create first three list when creating the board
            const ToDo = new List({
                listId: "L" + nanoid(),
                title: "To Do",
                index: 1,
                createdBy: email,
                boardId
            })
            const InProgress = new List({
                listId: "L" + nanoid(),
                title: "In Progress",
                index: 2,
                createdBy: email,
                boardId
            })
            const Completed = new List({
                listId: "L" + nanoid(),
                title: "Completed",
                index: 3,
                createdBy: email,
                boardId
            })

            await ToDo.save()
            await InProgress.save()
            await Completed.save()

            res.status(201).json({
                boardId: newboard.boardId,
                boardName: newboard.boardName,
                boardType: newboard.type
            })
        } catch (error) {
            console.log(error)
            res.status(500).send(error)
        }

        // newboard.save().then( () => {
        //     res.status(201).json({
        //         boardId: newboard.boardId,
        //         boardName: newboard.boardName,
        //         boardType: newboard.type
        //     })
        // }).catch(error => {
        //     console.log(error)
        //     res.status(500).send(error)
        // })
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
