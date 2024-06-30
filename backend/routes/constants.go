package routes

const userNameRegex = `^[a-zA-Z0-9 ]{5,32}$`
const accountNumberRegex = `^[a-f0-9]{8}-([a-f0-9]{4}-){3}[a-f0-9]{12}$`
const listIdRegex = `^[a-f0-9]{8}-([a-f0-9]{4}-){3}[a-f0-9]{12}$`

// TODO: allow more characters
const todoDescriptionRegex = `^[a-zA-Z0-9 ]{1,256}$`
