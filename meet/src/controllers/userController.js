const express = require('express')

const router = express.Router()
const parseJson = express.json({ extended: true })

// Import User model
const User = require('../models/user')

router
    .route('/new')
    .post([parseJson], async (req, res) => {
        try {
            const { member, projectid } = req.body
            const { name, fid, email, imgUrl, role } = member

            const projects = {
                projectId: projectid,
                role
            }

            const doc = await User.findOneAndUpdate({ email }, { name, fid, email, imgUrl, $push: { projects } }, { upsert: true, new: true })
            res.status(200).send("Content added")
        } catch (error) {
            console.log(error)
            res.status(500).send(error)
        }
    })

module.exports = router