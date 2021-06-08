const auth = require('../config/firebase')

async function fetchUser(uid) {
    try {
        const decodedToken = await auth.verifyIdToken(uid)
        return decodedToken
    } catch (error) {
        console.log(error)
        return null
    }
}

module.exports = fetchUser