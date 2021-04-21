const express = require('express')
const { customAlphabet } = require('nanoid');

const router = express.Router()
const parseJson = express.json({ extended: true })

const channelModel = require('../models/channelModel')
const userModel = require('../models/userModel')


const alphabet = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
const nanoid = customAlphabet(alphabet, 12);


router
    .route('/everyone/new')
    .post([parseJson], async (req, res) => {
        const member = {}
        const newChannel = new channelModel({
            projectid: req.body.projectid,
            channelName: '#everyone',
            channelid: nanoid(),

        })
        member.email = req.body.member.email
        member.isAdmin = false
        newChannel.members.push(member)
        newChannel.save()
            .then(doc => {
                console.log(doc)
                res.end()
            })
            .catch(error => {
                console.log(error)
                res.end()
            })
        const { name, fid, imgUrl, email } = req.body.member
        const newUser = new userModel({
            name,
            fid,
            imgUrl,
            email
        })
        newUser.role.push({
            projectid: req.body.projectid,
            role: req.body.member.role
        })
        newUser.save()
            .then(doc => {
                console.log(doc)
                res.end()
            })
            .catch(error => {
                console.log(error)
                res.end()
            })
    })

router
    .route('/new')
    .post([parseJson], async (req, res) => {
        const { projectid, channelName, members } = req.body
        members.push({
            email: req.decodedToken.email,
            isAdmin: true
        })
        const newChannel = new channelModel({
            projectid,
            channelName,
            channelid: nanoid(),
            members
        })
        newChannel.save()
            .then(doc => {
                res.status(201).json({
                    channelid: doc.channelid,
                    channelName: doc.channelName,
                    isAdmin: true
                })
            })
            .catch(error => {
                res.status(500).send(error)
            })
    })

module.exports = router