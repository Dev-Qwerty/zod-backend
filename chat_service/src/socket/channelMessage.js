const express = require('express')
const { nanoid } = require('nanoid')
const channelModel = require('../models/channelModel')
const userModel = require('../models/userModel')
const chatModel = require('../models/chatModel')
const VerifyUser = require('../middlewares/verifyuser')

const parseJson = express.json({ extended: true })

const channelMessage = async (projectSpace, socket, app) => {
    const projectid = projectSpace.name.slice(1).replace('/chat', "")
    const doc = await channelModel.find({ projectid, members: { $elemMatch: { email: socket.email } } }, 'channelid -_id')
    for (i = 0; i < doc.length; i++) {
        socket.join(doc[i].channelid)
    }

    app.post('/api/chat/:channelid/messages', [VerifyUser, parseJson], async (req, res) => {
        const channel = req.params.channelid
        const email = req.decodedToken.email
        const user = await userModel.findOne({ email }, 'name imgUrl -_id')
        const message = new chatModel({
            messageid: nanoid(),
            projectid,
            channelid: req.params.channelid,
            author: {
                email,
                name: user.name,
                imgUrl: user.imgUrl,
            },
            ts: Date.now(),
            content: req.body.content
        })
        await message.save()
        const respMessage = {
            projectid,
            channelid: message.channelid,
            content: message.content,
            author: message.author,
            messageid: message.messageid,
            ts: message.ts,
            edited: message.edited
        }

        projectSpace.to(channel).emit('newMessage', respMessage)
        res.send(respMessage)
    })

    app.delete('/api/chat/:channelid/messages/:messagets', [VerifyUser, parseJson], async (req, res) => {
        try {
            const channelid = req.params.channelid
            const ts = parseInt(req.params.messagets)
            const email = req.decodedToken.email
            const deletedMessage = await chatModel.findOneAndDelete({ channelid, ts, "author.email": email })
            if (deletedMessage) {
                projectSpace.to(channelid).emit("deleteMessage", { "channelid": deletedMessage.channelid, "ts": deletedMessage.ts })
                res.status(200).json({ "channelid": deletedMessage.channelid, "ts": deletedMessage.ts })
                return
            }
            res.status(400).send("Failed to delete message")
        } catch (error) {
            console.log("Delete Chat: ", error)
            res.status(500).send(error)
        }
    })

    app.put('/api/chat/:channelid/messages/:messagets/update', [VerifyUser, parseJson], async (req, res) => {
        try {
            const channelid = req.params.channelid
            const ts = parseInt(req.params.messagets)
            const email = req.decodedToken.email
            const updatedMessage = await chatModel.findOneAndUpdate({ channelid, ts, "author.email": email }, { content: req.body.content, edited: true }, { new: true }).lean()
            if (updatedMessage) {
                delete updatedMessage._id
                delete updatedMessage.__v
                projectSpace.to(channelid).emit("udpateMessage", updatedMessage)
                res.status(200).send(updatedMessage)
                return
            }
            res.status(400).send("Failed to updated message")
        } catch (error) {
            console.log(error)
            res.status(500).send(error)
        }
    })

}

module.exports = channelMessage