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
            newChannel.save()
            res.end()
        } catch (error) {
            console.log(error)
            res.end()
        }
    })

module.exports = router