const boardController = require('../controllers/boardController')
const itemController = require('../controllers/itemController')
const isUserVerified = require('../middlewares/verifyUser')

const router = (app) => {

    app.use('/api/board', isUserVerified, boardController)
    app.use('/api/item', isUserVerified, itemController)
    
}

module.exports = router
