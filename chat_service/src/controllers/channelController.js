const express = require('express')
const { customAlphabet } = require('nanoid');

const router = express.Router()
const parseJson = express.json({ extended: true })

const channelModel = require('../models/channelModel')


const alphabet = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
const nanoid = customAlphabet(alphabet, 12);


// router
//     .route('/everyone/new')
//     .post([parseJson], async (req, res) => {
//         const member = {}
//         const newChannel = new channelModel({
//             projectid: req.body.projectid,
//             channelName: '#everyone',
//             channelid: nanoid(),
//         })
//         member.email = req.body.member.email
//         member.isAdmin = false
//         newChannel.members.push(member)
//         newChannel.save()
//             .then(doc => {
//                 console.log(doc)
//                 res.end()
//             })
//             .catch(error => {
//                 console.log(error)
//                 res.end()
//             })
//         const { name, fid, imgUrl, email } = req.body.member
//         projectRole = []
//         projectRole.push({
//             projectid: req.body.projectid,
//             role: req.body.member.role
//         })
//         userModel.findOneAndUpdate({ email }, { name, fid, email, imgUrl, $push: { projectRole } }, { upsert: true, new: true })
//             .then(doc => {
//                 console.log(doc)
//                 res.end()
//             })
//             .catch(error => {
//                 console.log(error)
//                 res.end()
//             })
//     })

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


router
    .route('/delete')
    .delete([parseJson], async (req, res) => {
        try {
            const { projectid, channelid } = req.body
            const email = req.decodedToken.email
            const doc = await channelModel.findOne({ projectid, channelid, members: { $elemMatch: { email } } }, 'members.$ -_id')
            if (doc == null) res.status(401).send('Unauthorized')
            if (doc.members[0].isAdmin == true) {
                let doc = await channelModel.findOneAndDelete({ projectid, channelid })
                res.status(200).json({ 'channelid': doc.channelid, 'message': 'channel deleted' })
            } else {
                res.status(401).send('Unauthorized')
            }
        } catch (error) {
            console.log(`Delete channel: ${error}`)
            res.status(500).send('Internal Server Error')
        }
    })

router
    .route('/:projectid')
    .get(async (req, res) => {
        try {
            const projectid = req.params.projectid
            const email = req.decodedToken.email
            const channels = await channelModel.find({ projectid, members: { $elemMatch: { email } } }, 'channelid channelName -_id')
            if (channels == null) {
                res.status(400).send('Bad request')
                return
            }
            res.status(200).send(channels)
        } catch (error) {
            console.log(`Get channels: ${error}`)
            res.status(500).send(error)
        }
    })

router
    .route('/:projectid/:channelid')
    .get(async (req, res) => {
        const { projectid, channelid } = req.params
        const email = req.decodedToken.email
        const doc = await channelModel.findOne({ projectid, channelid, members: { $elemMatch: { email } } }, 'members.isAdmin.$ -_id')
        if (doc == null || doc.members[0].isAdmin == false) res.status(401).send('Unauthorized')
        if (doc.members[0].isAdmin == true) res.status(200).send(doc.members[0].isAdmin)
        // TODO: fetch messages
    })

router
    .route('/:projectid/:channelid/members')
    .get(async (req, res) => {
        try {
            const { projectid, channelid } = req.params
            const email = req.decodedToken.email
            let doc = await channelModel.findOne({ projectid, channelid, members: { $elemMatch: { email } } })
            if (doc == null) {
                res.status(400).send('Bad request')
                return
            }
            doc = await channelModel.aggregate([
                {
                    $match: {
                        projectid, channelid, members: { $elemMatch: { email } }
                    },
                },
                {
                    $lookup: {
                        from: 'users',
                        localField: 'members.email',
                        foreignField: 'email',
                        as: 'channelMembers'
                    }
                },
                {
                    $project: {
                        'channelMembers.name': 1,
                        'channelMembers.email': 1,
                        'channelMembers.imgUrl': 1
                    }
                }
            ])
            res.status(200).send(doc[0].channelMembers)
        } catch (error) {
            console.log(error)
            res.status(500).send('Internal Server Error')
        }

    })

module.exports = router