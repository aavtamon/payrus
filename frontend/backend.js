// PizTec Corporation, 2024. All Right Reserved

const auth_token_item = "auth_token"
const account_item = "account_token"
const auth_token_header = "Payrus-Auth-Token"

const BACKEND_CODE_SUCCESS = 200
const BACKEND_CODE_NOT_FOUND = 404
const BACKEND_CODE_ALREADY_EXISTS = 409
const BACKEND_CODE_SERVER_ERROR = 500

var Backend = {
    logIn: function(email, password, callback) {
        this._communicate("login", "POST", {email: email, password: password}, {
            success: function(account, status, request) {
                if (account == null) {
                    if (callback) {
                        callback(BACKEND_CODE_SERVER_ERROR, "Incorrect server response");
                    }
                } else {
                    window.sessionStorage.setItem(auth_token_item, request.getResponseHeader(auth_token_header))
                    window.sessionStorage.setItem(account_item, JSON.stringify(account))
                      
                    if (callback) {
                      callback(BACKEND_CODE_SUCCESS);
                    }
                }
            }.bind(this),
            error: function(errorText, status, request) {
                if (callback) {
                    callback(status, errorText);
                }
            }
        })
    },

    logOut: function(callback) {
        function removeSession() {
            window.sessionStorage.removeItem(auth_token_item)
            window.sessionStorage.removeItem(account_item)
        }

        this._communicate("logout", "POST", null, {
            success: function() {
                removeSession()
                
              if (callback) {
                callback(BACKEND_CODE_SUCCESS);
              }
            }.bind(this),
            error: function(errorText, status, request) {
                removeSession()
                
                if (callback) {
                    callback(status, errorText);
                }
            }
        })
    },

    isLogged: function() {
        return window.sessionStorage.getItem(auth_token_item) !== null && window.sessionStorage.getItem(account_item) !== null
    },

    createAccount: function(email, password, dob, callback) {
        this._communicate("create_account", "POST", {email: email, password: password, dob: dob}, {
            success: function() {
              if (callback) {
                callback(BACKEND_CODE_SUCCESS);
              }
            },
            error: function(errorText, status, request) {
              if (callback) {
                callback(status, errorText);
              }
            }
        })
    },

    changePassword: function(newPassword, currentPassword, callback) {
        this._communicate("change_account", "PUT", {password: newPassword, current_password: currentPassword}, {
            success: function() {
              if (callback) {
                callback(BACKEND_CODE_SUCCESS);
              }
            },
            error: function(errorText, status, request) {
              if (callback) {
                callback(status, errorText);
              }
            }
        })
    },

    orderCard: function(firstName, lastName, callback) {
        account = this.getAccount()
        if (account == null) {
            callback(BACKEND_CODE_NOT_FOUND);
            return
        }

        this._communicate("create_card", "POST", {account_id: account.id, first_name: firstName, last_name: lastName}, {
            success: function(account) {
                window.sessionStorage.setItem(account_item, JSON.stringify(account))

                if (callback) {
                    callback(BACKEND_CODE_SUCCESS);
                }
            },
            error: function(errorText, status, request) {
                if (callback) {
                    callback(status, errorText);
                }
            }
        })
    },

    getAccount: function() {
        account = window.sessionStorage.getItem(account_item)
        return account === null ? null : JSON.parse(account)
    },


    _communicate: function(resource, method, data, callback) {
        var request = new XMLHttpRequest();
        request.onreadystatechange = function() {
            if (request.readyState == request.DONE) {
                if (request.status >= 200 && request.status < 300) {
                    var text = request.responseText;
                    try {
                        text = JSON.parse(request.responseText);
                        callback.success(text, request.status, request);
                    } catch (e) {
                        callback.error(e, request.status, request);
                    }                    
                } else {
                    callback.error(request.responseText, request.status, request);
                }
            }
        }
    
        
        var url = "/api/" + resource;
    
        request.open(method, url, true);
        request.setRequestHeader("content-type", "application/json");
        request.setRequestHeader(auth_token_header, window.sessionStorage.getItem(auth_token_item));
    
        request.send(data != null ? JSON.stringify(data) : "");    
    }
}
