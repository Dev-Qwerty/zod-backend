const express = require('express')
const channelModel = require('../models/channelModel')
const userModel = require('../models/userModel')
const chatModel = require('../models/chatModel')

const parseJson = express.json({ extended: true })

const channelMessage = async (projectSpace, socket, app) => {
    const projectid = projectSpace.name.slice(1)
    const doc = await channelModel.find({ projectid, members: { $elemMatch: { email: socket.email } } }, 'channelid -_id')
    for (i = 0; i < doc.length; i++) {
        socket.join(doc[i].channelid)
    }

    app.post('/api/chat/:channelid/messages', [parseJson], async (req, res) => {
        const channel = req.params.channelid
        const user = await userModel.findOne({ email: socket.email }, 'name imgUrl -_id')
        const message = new chatModel({
            projectid,
            channelid: req.params.channelid,
            author: {
                email: socket.email,
                name: user.name,
                imgUrl: user.imgUrl,
            },
            content: req.body.content
        })
        await message.save()
        const respMessage = {
            id: message._id,
            projectid,
            channelid: message.channelid,
            content: message.content,
            author: message.author,
            date: message.date
        }

        projectSpace.to(channel).emit('new message', respMessage)
        res.send(respMessage)
    })

}

module.exports = channelMessage