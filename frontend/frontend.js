// PizTec Corporation, 2024. All Right Reserved

const auth_token_item = "auth_token"

function createAccount(email, password, dob, callback) {
    window.sessionStorage.setItem(auth_token_item, "test")

    callback(null)
}

function logIn(email, password, callback) {
    window.sessionStorage.setItem(auth_token_item, "test")
    callback(null)
}

function logOut(callback) {
    window.sessionStorage.removeItem(auth_token_item)

    callback(null)
}

function getToken() {
    return window.sessionStorage.getItem(auth_token_item)
}