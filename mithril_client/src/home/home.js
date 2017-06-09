var m = require("mithril");

module.exports = {
    view: function() {
        return m(".bg-light-green", [
            m(".cover.ba.slant-cover", {style: "background: white url('/assets/img/business-heights.jpeg') center no-repeat"},
                m(".pa2.h5", {style: "background-color: rgba(0, 0, 0, .5)"},
                    m("nav.pv2.ph7.cf",
                        m(".cf", [
                            m("h2.mv0.fl.white", "DEE WEB CHAT"),
                            m("a.white.no-underline.fr.mv1[href='https://github.com/spankie/web-chat'][target='_blank']",
                                m("img", {height: "25", width: "25", src: "/assets/img/GitHub-Mark-32px.png", alt:"GitHub", class:"dib v-top"})
                            )
                        ])
                    )
                )
            ),
            m(".ph1.ph5-ns.tc.relative", {style:"bottom: 12rem"}, [
                m("p.f2.white.mv1.w6", {style: "margin-left: auto; margin-right: auto"}, "LOGIN"),
                m(".dib.pa3.ba.bw2.b--green.br1.shadow-5.w6", {style: "background: white url('/assets/img/grid-bg.png') top left repeat"}, [
                    m("p.pa2.br1.bg-light-red.shadow-2.white.f6.tl.mt0.mb1", "ERROR || MESSAGE"),
                    m("form", [
                        m(".mb2.tl", [m("label.db.mb2.gray", "Username"), m("input.pa2.ba.b--green.br1.w-100", {type: "text", name:"username", placeholder: "Username"})]),
                        m(".mb2.tl", [m("label.db.mb2.gray", "Password"), m("input.pa2.ba.b--green.br1.w-100", {type: "password", name: "passwd", placeholder: "Password"})]),
                        m(".mv3.tl", [
                            m("p.tc.gray", "Login or Sign up"),
                            m("input.pv2.ph3.ba.bg-white.green.hover-white.br1.shadow-5.pointer", {type:"submit", name:"login", value:"LOGIN"}),
                            m("button.fr.pv2.ph3.ba.bg-white.green.hover-bg-green.hover-white.br1.shadow-5.pointer", "Sign Up")
                        ])
                    ])
                ])
            ])
        ])
    }
}