const mongoose = require('mongoose')

const userSchema = new mongoose.Schema({
    name: String,
    id: String
})

const cardSchema = new mongoose.Schema({
    cardId: String,
    cardTitle: String,
})

const itemsSchema = new mongoose.Schema({
    title: String,
    description: String,
    createdBy: userSchema,
    assignedMember: [userSchema],
    duedate: String,
    boardId: String,
    cardName: cardSchema
})