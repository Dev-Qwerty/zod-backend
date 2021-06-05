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
            // Check if the user is a member of project
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

            // Fetch lists and cards in the board
            lists = await List.aggregate([
                { $match: { boardId } },
                {
                    $lookup: {
                        from: "cards",
                        localField: "listId",
                        foreignField: "listId",
                        as: "cards"
                    }
                },
                {
                    $project: {
                        __v: 0,
                        _id: 0,
                        "cards.__v": 0,
                        "cards._id": 0,
                    }
                }
            ])

            // Fetch members in the board
            doc = await Board.aggregate([
                { $match: { boardId } },
                {
                    $lookup: {
                        from: "users",
                        localField: "members.email",
                        foreignField: "email",
                        as: "members"
                    }
                },
                {
                    $project: {
                        _id: 0,
                        "members.name": 1,
                        "members.email": 1,
                        "members.imgUrl": 1
                    }
                }
            ])

            const response = {
                members: doc[0].members,
                lists
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
            doc = await User.find({
                projects: {
                    $elemMatch: { projectId }
                }
            }, 'email')

            i = 0
            while (i < members.length) {
                j = 0
                while (j < doc.length) {
                    if (members[i].email == doc[j].email) {
                        break
                    }

                    if (j == doc.length - 1) {
                        members.splice(i, 1)
                        i--
                    }

                    j++
                }
                i++
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

            // Fetch all members in the project
            doc = await User.find({
                projects: {
                    $elemMatch: { projectId }
                }
            }, 'email')

            // Check if the members in req are in the project
            i = 0
            while (i < members.length) {
                j = 0
                while (j < doc.length) {
                    if (members[i].email == doc[j].email) {
                        break
                    }

                    if (j == doc.length - 1) {
                        members.splice(i, 1)
                        i--
                    }
                    j++
                }
                i++
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
                createdBy: email,
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
            // Check if the user is a member of project
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
