const admin = require('firebase-admin')

const sa = process.env.FIREBASE_SA_KEY

admin.initializeApp({
    credential: admin.credential.cert(sa)
})

const auth = admin.auth()

module.exports = auth