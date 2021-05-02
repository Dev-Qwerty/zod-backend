// User model
const User = require('../models/user')

// Add user details to db
// When the user creates or accept a project
const newUser = (value) => {
    const { name, fid, email, imgUrl, } = value.member
    const projects = value.projects
    User.findOneAndUpdate({ email }, { name, fid, email, imgUrl, $push: { projects } }, { upsert: true, new: true })
        .then(doc => {
            console.log(doc)
        }).catch(error => {
            console.log(error)
        })
}

module.exports = newUser
