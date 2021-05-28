const express = require('express')

const router = express.Router()

const chatModel = require('../models/chatModel')
const channelModel = require('../models/channelModel')

router
    .route('/:channelid')
    .get(async (req, res) => {
        const latest = parseInt(req.query.latest)
        const channelid = req.params.channelid
        const email = req.decodedToken.email
        let doc = await channelModel.findOne({ channelid, members: { $elemMatch: { email } } })
        if (doc == null) {
            res.status(400).send("Bad Request")
            return
        }
        let messages = await chatModel.find({ channelid, ts: { $lt: latest } }).sort({ ts: -1 }).limit(20)
        messages.reverse()
        res.status(200).send(messages)
    })


module.exports = router