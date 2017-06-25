"use strict";

var m     = require("mithril");
var Home  = require("./home/home");
var User  = require("./user/user");

var root = document.body;
// m.render(document.body, m("h1.bg-red", "Hello there!! This is web-chat with mithril"));
m.route(root, "/", {
    "/": {
        view: () => {
            // return m("h1.bg-red", "Hello there!! This is web-chat with mithril")
            return m(Home);
        }
    },
    "/:username": {
        view: () => {
            return m(User);
        }
    }
})