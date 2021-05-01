const { customAlphabet } = require('nanoid');
const alphabet = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
const nanoid = customAlphabet(alphabet, 12);

const channelModel = require('../models/channelModel')
const userModel = require('../models/userModel')

const CreateChannelEveryone = async (message) => {
    const member = {}
    const newChannel = new channelModel({
        projectid: message.projectid,
        channelName: '#everyone',
        channelid: nanoid(),

    })
    member.email = message.member.email
    member.isAdmin = false
    newChannel.members.push(member)
    newChannel.save()
        .then(doc => {
            console.log(doc)
        })
        .catch(error => {
            console.log(error)
        })
    const { name, fid, imgUrl, email } = message.member
    projectRole = []
    projectRole.push({
        projectid: message.projectid,
        role: message.member.role
    })
    userModel.findOneAndUpdate({ email }, { name, fid, email, imgUrl, $push: { projectRole } }, { upsert: true, new: true })
        .then(doc => {
            console.log(doc)
        })
        .catch(error => {
            console.log(error)
        })
}

module.exports = {
    CreateChannelEveryone
}