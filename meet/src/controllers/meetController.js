const express = require('express')
const { customAlphabet } = require('nanoid')

const router = express.Router()

// Import Models
const Meet = require('../models/meet')
const User = require('../models/user')

const parseJson = express.json({ extended: true })

// Create 12 digit alphanumeric code for boardId and chardId
const keyStore = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const nanoid = customAlphabet(keyStore, 12)

router
    .route('/new')
    .post([parseJson], async (req, res) => {
        try {
            const { meetName, projectId, members, createdBy } = req.body

            const meetId = nanoid()
            const meetUrl = `${process.env.BASE_URL}/${meetId}`
            console.log(meetUrl)

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
                    if (members[i] == doc[j].email) {
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

            members.push(createdBy)

            let newMeet = new Meet({
                meetId,
                meetName,
                meetUrl,
                createdBy,
                projectId,
                members
            })

            await newMeet.save()

            res.status(201).json({ link: meetUrl })
        } catch (error) {
            console.log(error)
            res.status(500).send(error)
        }

    })

module.exports = router