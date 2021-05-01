const boardController = require('../controllers/boardController')
const userController = require('../controllers/userController')
const isUserVerified = require('../middlewares/verifyUser')

const router = (app) => {
    
    app.use('/api/everyone', userController.router)
    app.use('/api/board', isUserVerified, boardController)
    
}

module.exports = router
