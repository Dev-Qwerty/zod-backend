const channelModel = require('../models/channelModel')

const channelMessage = async (projectSpace, socket) => {
    const projectid = projectSpace.name.slice(1)
    const doc = await channelModel.find({ projectid, members: { $elemMatch: { email: socket.email } } }, 'channelid -_id')
    for (i = 0; i < doc.length; i++) {
        socket.join(doc[i].channelid)
    }
    console.log(socket.rooms);
}

module.exports = channelMessage