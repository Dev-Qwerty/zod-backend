const boardController = require('../controllers/boardController')
const cardController = require('../socket/boardChannel')
const listController = require('../controllers/listController')
const isUserVerified = require('../middlewares/verifyUser')

const router = (app) => {

    app.use('/api/board', isUserVerified, boardController)
    
}

module.exports = router
