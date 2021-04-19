const boardController = require('../controllers/boardController')

const router = (app) => {
    
    app.use('/api/board', boardController)
    
}

module.exports = router