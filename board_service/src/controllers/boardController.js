const express = require('express')
const { customAlphabet } = require('nanoid')

// Import models
const Board = require('../models/board')
const User = require('../models/user')
const List = require('../models/list')
const Card = require('../models/card')

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

// @route GET /api/board/lists/:boardid
// @desc Fetch the list and cards in the board
router
    .route('/lists/:boardId')
    .get(async (req, res) => {
        try {
            const boardId = req.params.boardId

            // Fetch email from Firebase decodedToken
            const email = req.decodedToken.email

            // Check if the user is the member of the board
            let doc = await Board.findOne({
                boardId, members: {
                    $elemMatch: { email }
                }
            })

            if (doc == null) {
                res.status(401).send("Unauthorized user")
            }

            const list = await List.find({ boardId })
            const card = await Card.find({ boardId })

            const response = {
                lists: list,
                cards: card
            }

            res.status(200).send(response)

        } catch (error) {
            console.log(error)
            res.status(500).send(error)
        }
    })

// @route POST /api/board/member/new
// @desc Add new members to the board
router
    .route('/member/new')
    .post([parseJson], async (req, res) => {
        try {
            const { boardId, projectId, members } = req.body

            const email = req.decodedToken.email

            // Check if the user is inside the project(email owner)
            let doc = await User.findOne({
                email, projects: {
                    $elemMatch: { projectId }
                }
            })

            if (doc == null) {
                res.status(401).send("Unauthorized user")
                return
            }

            // Check if the members are in the project
            for (i = 0; i < members.length; i++) {
                const email = members[i].email
                let doc = await User.findOne({
                    email, projects: {
                        $elemMatch: { projectId }
                    }
                })

                if (doc == null) {
                    members.splice(i, 1)
                }
            }

            await Board.findOneAndUpdate({ boardId }, { $push: { members } })

            res.status(200).send(members)

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

            // Check if the user is inside the project(email owner)
            let doc = await User.findOne({
                email, projects: {
                    $elemMatch: { projectId }
                }
            })

            if (doc == null) {
                res.status(401).send("Unauthorized user")
                return
            }

            // Check if members are inside the project(project members)
            for (i = 0; i < members.length; i++) {
                const email = members[i].email
                let doc = await User.findOne({
                    email, projects: {
                        $elemMatch: { projectId }
                    }
                })

                if (doc == null) {
                    delete members[i]
                }
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

            res.status(201).json({
                boardId: newboard.boardId,
                boardName: newboard.boardName,
                boardType: newboard.type
            })
        } catch (error) {
            console.log(error)
            res.status(500).send(error)
        }
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
            let doc = await Board.findOne({
                boardId, projectId, members: {
                    $elemMatch: { email }
                }
            }, 'members.$ -_id')

            if (doc == null) {
                res.status(401).send("Unauthorized user")
                return
            }

            if (doc.members[0].isAdmin == true) {
                let doc = await Board.findOneAndDelete({ boardId, projectId })
                await List.deleteMany({ boardId })
                await Card.deleteMany({ boardId })
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
