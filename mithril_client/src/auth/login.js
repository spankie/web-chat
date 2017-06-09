var m = require("mithril");
var data = require("../models/data");

module.exports = () => {
    // get the username and password...
    return m.request({
        method: "POST",
        url: "/api/login",
        data: {username: data.myUsername, password: data.myPassword},
        withCredentials: true
    })
    .then(function(result) {
        // check if the result is an error or contains a cookie...
        console.log("login result:" + result);
    })
}