const express = require('express')

const router = express.Router()
const parseJson = express.json({ extended: true })

// User model
const User = require('../models/user')

// Add user details to db
// When the user creates or accept a project
const newUser = (value) => {
    const { name, fid, email, imgUrl, } = value.member
    const projects = value.projects
    User.findOneAndUpdate({ email }, { name, fid, email, imgUrl, $push: { projects } }, { upsert: true, new: true })
        .then(doc => {
            console.log(doc)
        }).catch(error => {
            console.log(error)
        })
}

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

            console.log(email)

            const doc = await User.findOneAndUpdate({ email }, { name, fid, email, imgUrl, $push: { projects } }, { upsert: true, new: true })
            res.status(200).send("Content added")
        } catch (error) {
            console.log(error)
            res.status(500).send(error)
        }
    })

module.exports = { newUser, router }
