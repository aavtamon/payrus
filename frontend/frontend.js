// PizTec Corporation, 2024. All Right Reserved

const auth_token_item = "auth_token"

accountConfiguration = {
    conversionRate: 100
}

function createAccount(email, password, dob, callback) {
    window.sessionStorage.setItem(auth_token_item, "test")

    callback(null)
}

function orderCard(firstName, lastName, callback) {
    accountConfiguration.card = {
        firstName: firstName,
        lastName: lastName,
        expiration: "2/26",
        balance: 100,
    }
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

function getConfiguration() {
    return accountConfiguration
}


function validateInputField(elementId, statusElementId, validator) {
    element = document.getElementById(elementId)
    statusElement = document.getElementById(statusElementId)

    isValid = false

    if (element.value == "") {
      statusElement.classList.remove("bad")
      statusElement.classList.remove("good")
    } else {
      if (validator(element.value)) {
        statusElement.classList.remove("bad")
        statusElement.classList.add("good")

        isValid = true
      } else {
        statusElement.classList.remove("good")
        statusElement.classList.add("bad")
      }
    }

    return isValid
}

function setElementEnabled(elementName, isEnabled) {
    element = document.getElementById(elementName)
    if (isEnabled) {
        element.removeAttribute("disabled")
    } else {
        element.setAttribute("disabled", "")
    }
}

function setElementValid(elementName, isValid) {
    element = document.getElementById(elementName)
    if (isValid) {
        element.classList.remove("invalid")
    } else {
        element.classList.add("invalid")
    }
}
