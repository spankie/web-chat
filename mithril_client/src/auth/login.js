var m = require("mithril");
var data = require("../models/data");

module.exports = () => {
    if(data.myUsername == "" || data.myPassword == ""){
        data.authError = "Username or password cannot be empty.";
        return;
    }
    // get the username and password...
    console.log("login");
    return m.request({
        method: "POST",
        url: "/api/login",
        data: {username: data.myUsername, password: data.myPassword},
        withCredentials: true
    })
    .then(function(response) {
        // check if the result is an error or contains a cookie...
        console.log("login response:", response);
        if(response.status == "ok") {
            // clear error if there is any showing
            data.authError = "";
            data.authMessage = "logged in";
            // TODO: add expiry date...
            document.cookie = "deewebchat=" + response.cookie + ";path=/";
            // TODO: redirect to /web/myUsername
            m.route.set("/" + data.myUsername);
        } else if (response.status == "error") {
            // clear message if there is any showing
            data.authMessage = "";
            // set the error to display...
            data.authError = response.error;
        } else {
            // incase of unforseen circumstances
            console.log("You are not logged in...");
        }
    })
}