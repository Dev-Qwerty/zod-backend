const express = require('express')

// User model
const user = require('../models/user')

const router = express.Router()
const parseJson = express.json({ extended: true })

// @route POST /api/everyone
// @desc Add users when a new project is created
router
    .route('/')
    .post( [parseJson], (req, res) => {
        const { projectId, members } = req.body
        
        newuser = new user({
            fid: members.fid,
            name: members.name,
            email: members.email,
            imgUrl: members.imgUrl,
        })
        
        newuser.role.push({
            projectId: projectId,
            role: members.role
        })

        newuser.save().then( () => {
            res.status(201).end()
        })
    })

module.exports = router