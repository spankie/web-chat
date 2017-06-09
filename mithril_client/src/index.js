var m = require("mithril");
var Home = require("./home/home");
var User = require("./user/user");

// m.render(document.body, m("h1.bg-red", "Hello there!! This is web-chat with mithril"));
m.route(document.body, "/", {
    "/": {
        render: function() {
            // return m("h1.bg-red", "Hello there!! This is web-chat with mithril")
            return m(Home);
        }
    },
    "/:username": {
        render: function() {
            return m(User);
        }
    }
})