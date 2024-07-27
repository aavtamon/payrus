// PizTec Corporation, 2024. All Right Reserved

const auth_token_item = "auth_token"

var Backend = {
    accountConfiguration: {
        conversionRate: 100
    },

    createAccount: function(email, password, dob, callback) {
        window.sessionStorage.setItem(auth_token_item, "test")

        callback(null)
    },

    orderCard: function(firstName, lastName, callback) {
        accountConfiguration.card = {
            firstName: firstName.toUpperCase(),
            lastName: lastName.toUpperCase(),
            expiration: "2/26",
            balance: 100,
        }
        callback(null)
    },

    logIn: function(email, password, callback) {
        window.sessionStorage.setItem(auth_token_item, "test")
        callback(null)
    },

    logOut: function(callback) {
        window.sessionStorage.removeItem(auth_token_item)

        callback(null)
    },

    getToken: function() {
        return window.sessionStorage.getItem(auth_token_item)
    },

    getConfiguration: function() {
        return accountConfiguration
    }
}
