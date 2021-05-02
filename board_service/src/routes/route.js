const boardController = require('../controllers/boardController')
const userController = require('../controllers/userController')
const itemController = require('../controllers/itemController')
const isUserVerified = require('../middlewares/verifyUser')

const router = (app) => {
    
    app.use('/api/everyone', userController.router)
    app.use('/api/board', isUserVerified, boardController)
    app.use('/api/item', isUserVerified, itemController)
    
}

module.exports = router
