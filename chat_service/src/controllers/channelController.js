const express = require('express')
const { customAlphabet } = require('nanoid');

const router = express.Router()
const parseJson = express.json({ extended: true })

const channelModel = require('../models/channelModel')
const userModel = require('../models/userModel')


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


// @route POST /api/channel/new
// @desc Create new channel
router
    .route('/new')
    .post([parseJson], async (req, res) => {
        const { projectid, channelName, description, members } = req.body
        const checkProject = await channelModel.findOne({ projectid, channelName: 'everyone' })
        if (checkProject == null) {
            res.status(400).send('Bad request')
            return
        }
        const channel = await channelModel.findOne({ projectid, channelName })
        if (channel != null) {
            res.status(400).send('Channel Name already exists')
            return
        }
        members.push({
            email: req.decodedToken.email,
            isAdmin: true
        })
        const newChannel = new channelModel({
            projectid,
            channelName,
            channelid: nanoid(),
            description,
            members
        })
        newChannel.save()
            .then(doc => {
                res.status(201).json({
                    channelid: doc.channelid,
                    channelName: doc.channelName,
                    description: doc.description,
                    isAdmin: true
                })
            })
            .catch(error => {
                res.status(500).send(error)
            })
    })


// @route DELETE /api/channel/delete
// @desc delete channel
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


// @route GET /api/channel/:projectid
// @desc fetch all channels in a project
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



// @route GET /api/channel/:projectid/:channelid
// @desc fetch channel details
router
    .route('/:projectid/:channelid')
    .get(async (req, res) => {
        try {
            const { projectid, channelid } = req.params
            const email = req.decodedToken.email
            const doc = await channelModel.findOne({ projectid, channelid, members: { $elemMatch: { email } } }, 'members.isAdmin.$ -_id')
            if (doc == null || doc.members[0].isAdmin == false) res.status(401).send('Unauthorized')
            if (doc.members[0].isAdmin == true) res.status(200).json({ "isAdmin": doc.members[0].isAdmin })
            // TODO: fetch messages
        } catch (error) {
            console.log(`Get Channel: ${error}`)
        }

    })


// @route GET /api/channel/:projectid/:channelid/members
// @desc fetch channel members
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

    // @route POST /api/channel/:projectid/:channelid
    // @desc add new member to a channel
    .post([parseJson], async (req, res) => {
        const projectid = req.params.projectid
        const channelid = req.params.channelid
        const members = req.body.members
        const email = req.decodedToken.email

        const checkUserisAdmin = await channelModel.findOne({ projectid, channelid, members: { $elemMatch: { email, isAdmin: 'true' } } })
        if (checkUserisAdmin != null) {
            for (i = 0; i < members.length; i++) {
                const email = members[i].email
                const projectMember = await channelModel.findOne({ projectid, channelName: 'everyone', members: { $elemMatch: { email } } })
                if (projectMember != null) {
                    await channelModel.updateOne({ projectid, channelid }, { $push: { members: { email: email, isAdmin: "false" } } })
                }
            }
            const channelMembers = await channelModel.findOne({ projectid, channelid }, 'members -_id')
            res.status(200).send(channelMembers)
            return
        }
        res.status(401).send("Unauthorized")
    })

router
    .route('/:projectid/:channelid/fetchmembers')
    .get(async (req, res) => {
        const projectid = req.params.projectid
        const channelid = req.params.channelid
        let projectMembers = await userModel.find({ projectRole: { $elemMatch: { projectid } } }, 'email name imgUrl -_id')
        let channelMembers = await channelModel.findOne({ projectid, channelid }, 'members -_id')
        while (channelMembers.members.length > 0) {
            j = 0
            while (j < projectMembers.length) {
                if (projectMembers[j].email == channelMembers.members[0].email) {
                    projectMembers.splice(j, 1)
                    channelMembers.members.splice(0, 1)
                    break
                }
                j++
            }
        }
        res.send(projectMembers)
    })


module.exports = router