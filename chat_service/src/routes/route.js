const channelController = require('../controllers/channelController')
const chatController = require('../controllers/chatController')

const initializeRoutes = (app) => {
    app.use('/api/channel', channelController)
    app.use('/api/messages', chatController)
}

module.exports = initializeRoutes