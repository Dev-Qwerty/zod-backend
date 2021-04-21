const boardController = require('../controllers/boardController')
const isUserVerified = require('../middlewares/verifyUser')

const router = (app) => {
    
    app.use('/api/board', isUserVerified, boardController)
    
}

module.exports = router