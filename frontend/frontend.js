// PizTec Corporation, 2024. All Right Reserved

var Frontend = {
    emailValidator: function(value) {
        return value.toLowerCase().match(
            /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|.(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
        )
    },

    passwordValidator: function(value) {
        return value.toLowerCase().match(
            /^([0-9]|[a-z]){10,}$/
        )
    },

    validateInputField: function(elementId, validator) {
        element = document.getElementById(elementId)

        if (element.value == "") {
            element.classList.remove("valid")
            element.classList.remove("invalid")

            return false
        }

        if (validator(element.value)) {
            element.classList.remove("invalid")
            element.classList.add("valid")

            return true
        } else {
            element.classList.remove("valid")
            element.classList.add("invalid")

            return false
        }
    },


    /*
    validateInputFieldWithStatus: function(elementId, statusElementId, validator) {
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
    },
    */

    setElementEnabled: function(elementName, isEnabled) {
        element = document.getElementById(elementName)
        if (isEnabled) {
            element.removeAttribute("disabled")
        } else {
            element.setAttribute("disabled", "")
        }
    }
}

