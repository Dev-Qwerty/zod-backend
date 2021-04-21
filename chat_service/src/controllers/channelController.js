const express = require('express')
const { customAlphabet } = require('nanoid');

const router = express.Router()
const parseJson = express.json({ extended: true })

const channelModel = require('../models/channelModel')


const alphabet = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
const nanoid = customAlphabet(alphabet, 12);


router
    .route('/everyone/new')
    .post([parseJson], async (req, res) => {
        try {
            let newChannel = new channelModel;
            newChannel.projectid = req.body.projectid
            newChannel.members = req.body.members
            newChannel.channelName = '#everyone'
            newChannel.channelid = nanoid()
            newChannel.members[0].channelAdmin = false
            console.log(newChannel)
            newChannel.save()
            res.end()
        } catch (error) {
            console.log(error)
            res.end()
        }
    })

router
    .route('/new')
    .post([parseJson], async (req, res) => {
        try {
            const { projectid, channelName } = req.body
            let newChannel = new channelModel({
                projectid,
                channelName,
                channelid: nanoid()
            })
            let memberDetails = await channelModel.find({
                projectid,
                channelName: '#everyone',
                members: {
                    $elemMatch: {
                        fid: req.decodedToken.uid
                    }
                }
            },
                'members.$ -_id'
            )
            memberDetails[0].members[0].channelAdmin = true
            newChannel.members.push(memberDetails[0].members[0])
            for (let i = 0; i < req.body.members.length; i++) {
                let memberDetails = await channelModel.find({
                    projectid,
                    channelName: '#everyone',
                    members: {
                        $elemMatch: {
                            memberid: req.body.members[i].memberid
                        }
                    }
                },
                    'members.$ -_id'
                )
                newChannel.members.push(memberDetails[0].members[0])
            }
            newChannel.save()
            res.status(201).send('channel created')
        } catch (error) {
            console.log(error)
            res.status(500).send(error)
        }
    })

module.exports = router