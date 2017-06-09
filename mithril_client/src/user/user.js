var m = require("mithril");
var data = require("../models/data");

module.exports = {
    oninit: () => {
        // first check if the user is logged in...incase of trying to load the page without without loggin in
        // this should load friends from the server...
    },
    view: function(vnode) {
        // return m(".pa2.bg-black.white.tc.f3", "Hi " + vnode.attrs.username);
        return m(".body", [
            m(".shadow-4", m("nav", m(".bg-green.pv3.ph3",
                m("a.no-underline.dib.pa1", {href: "/", oncreate: m.route.link}, m("h2.mv0.white", "DEE WEB CHAT"))))
            ),

            m(".pa2", [
                m(".fl.w-25.pa2", m(".b--green.bg-white.pa1.shadow-4", [
                        m("div", [
                            m("h2.tc", "Users"),
                            m("hr.light"),
                            m(".tr", [
                                m("p.red.tc.mv1", "Error"),
                                m("input.pa2.w-100.mv1", {
                                    type:"text",
                                    oninput: m.withAttr("value", function(value) { data.friendName = value; }),
                                }),
                                m(".popup.relative.dib.pointer", m("div#myPopup.popuptext.white.tl.br3.pa2.absolute", [
                                    m("br-100.w3.h3.ba.b--gray.overflow-hidden.dib.v-mid", m("img", {src: "/assets/img/avatar.jpg", alt: "f.username"})),
                                    m(".dib.v-mid", [
                                        m("p.mv1.f4", "Friend.Username"),
                                        m("p.mv1.f6", "Status: " + "friend.id")
                                    ])
                                ])),
                                m("button.pa2.w-100.bg-green.ba.b--green.white.h2", [
                                    m("span", "Add"),
                                    m(".loader", "Loading") // ng-show(addLoader)
                                ])
                            ]),
                            m("hr.light")
                        ]),
                        m("div", {style:"overflow:auto;max-height:400px"}, data.friends.map(function(friend){
                            return m(".bb.pa2.b--light-gray", m(".dim.pointer", {id: friend.id}, [
                                m(".br-100.w3.h3.ba.b--gray.overflow-hidden.dib.v-mid", m("img", {src: "/assets/img/avatar.jpg", alt: friend.username})),
                                m(".dib.v-mid", [
                                    m("p.mv1.f4.gray", friend.username),
                                    m("p.mv1.f6.", "status: " + friend.id)
                                ])
                            ]))
                        })) 
                    ])
                ),
                m(".fl.w-50.pa2",),
                m(".fl.w-25.pa2",)
            ])
        ]);
    }
}
