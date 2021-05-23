const boardController = require('../controllers/boardController')
const userController = require('../controllers/userController')
const isUserVerified = require('../middlewares/verifyUser')

const router = (app) => {

    app.use('/api/board', isUserVerified, boardController)
    app.use('/api/user', userController.router)
    
}

module.exports = router
