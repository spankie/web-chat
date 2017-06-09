var m = require("mithril");
var data = require("../models/data");

module.exports = () => {
    if(data.myUsername == "" || data.myPassword == ""){
        data.authError = "Username or password cannot be empty.";
        return;
    }
    // signup the user here...
    console.log("signup");
    return m.request({
        method: "POST",
        url: "/api/signup",
        data: {username: data.myUsername, password: data.myPassword},
        withCredentials: true
    })
    .then(function(response) {
        // check if the response contains error(s)
        console.log("signup response: ", response);
        if(response.status == "ok"){
            // signup successful
            // clear error if there is any showing
            data.authError = "";
            data.authMessage = "signed up";
            // TODO: add expiry date...
            document.cookie = "deewebchat=" + response.cookie + ";path=/";
            // redirect to /web/myUsername
            m.route.set("/" + data.myUsername);
        } else if (response.status == "error") {
            // clear message if there is any showing
            data.authMessage = "";
            // set the error to display...
            data.authError = response.error;
        } else {
            data.authError = "";
            data.authMessage = "";
            // incase something else is returned
            console.log("You are not logged in.");
        }
    })
}