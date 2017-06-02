userarea = angular.module("userarea", ['ngCookies', 'ngWebSocket']);

userarea.controller("user", function($scope, $location, $cookies, $http, $window, $websocket) {
    var chatBox = document.getElementById("chatbox");
    chatBox.scrollTop = chatBox.scrollHeight;
    
    urlsplit = $location.absUrl().split("/");
    username = urlsplit[urlsplit.length - 1];
    $scope.username = username;
    $scope.friends = [];
    $scope.chat = {}
    // Fetch friends. (Incase of page reloads...)
    $http.post("/api/get/friends", {}).then(function(response){
        // success response callback
        data = response.data;
        if (data.hasOwnProperty('status') && data.status == 'error') {
            // You have no friends
            // tell the user
            console.log("no friends:", data, " type:", typeof data);
            return;
        } else if (data.length > 0) {
            $scope.friends = data;
            data.forEach(function(fr) {
                // each array of message is going to contain objects of messages...
                $scope.chat[fr.username] = [];
            }, this);
            console.log("Got Some friends");
            return;
        }
        console.log("neither no friends nor friends. Data: ", data, " type: ", typeof data);
    }, function(response) {
        // error response callback
        console.log("Error response: ". response)
    });

    var appendChat = (msg) => {
        $scope.$evalAsync(() => {
            if ($scope.context != null) {
                $scope.chat[$scope.context.username].push(msg);
            }
        })
    }

    // WEBSOCKET
    if(window["WebSocket"]) {
        console.log("WebSocket is available");
        // initiate the connection
        var conn = new WebSocket("ws://" + document.location.host + "/api/chat");
        conn.onclose = function (event) {
            console.log("connection closed.")
        }
        conn.onmessage = function (event) {
            var m = event.data;
            console.log("websocket message: " + m);
            if (m == "HI") {
                return;
            }
            
            var messageJSON = JSON.parse(m);
            messageJSON.Message = messageJSON.Message.replace("\n", "<br>");
            appendChat(messageJSON)
        }
    } else {
        console.log("NO WEBSOCKET...");
    }

    // $scope.friends = [{id: 60, username: "Silvia"}];
    // if a user search returns a result...
    $scope.foundFriend = false;
    $scope.friendName = "";
    $scope.friend = null;
    $scope.addLoader = false;
    $scope.context = null;
    $scope.message = "";

    $scope.sendMessage = () => {
        if (!conn) {
            console.log("No connection");
            $scope.errorMessage = "Connection not available";
            return
        }
        if ($scope.context === null || $scope.context.id == 'undefined' || isNaN($scope.context.id)) {
            console.log("COntext :", $scope.context);
            // console.log("typeof id", typeof $scope.context.id);
            return;
        }

        // check if message is empty:
        if ($scope.message == "") {
            $scope.errorMessage = "Empty message.";
            return;
        }
        // var theMessage = $scope.context.id + "\n" + $scope.message;
        // var mm = "2\nHello There";
        var now = new Date();
        now = now.getDate() + "/" + now.getMonth() + "/" + now.getFullYear();
        var theMessage = {
            Sender: $scope.username,
            Datetime: now,
            Message: $scope.message,
            Recepient: $scope.context.id
        }

        var messageJSON = JSON.stringify(theMessage);
        conn.send(messageJSON);
        console.log("Sent: ", messageJSON);
        // add the message to the chat history of this partiular context user
        // console.log("$Scope.chat: ", $scope.chat, " username: ", $scope.context.username);
        theMessage.Message = theMessage.Message.replace("\n", "<br>");
        appendChat(theMessage);
        $scope.message = "";
    }

    $scope.changeContext = (me) => {
        console.log(me.f.id)
        $scope.context = {
            id       : me.f.id,
            username : me.f.username
        }
        /*
            Message : {
                sender: me/friend
                datetime: new Date()
                message: "hello sire"
            }
        */
    }

    $scope.searchFriend = () => {
        console.log("FRIEND: ", $scope.friend)
        // search for a friend from the server.
        // /api/search/friend
        $http.post("/api/search/friend", {username: $scope.friendName}).then(function(response){
            // SUCCESS CALLBACK
            data = response.data;
            // console.log(data);
            if (data.hasOwnProperty('status') && data.status == 'error') {
                // Could not get the user...
                // Display the error here
                $scope.foundFriend = false;
                return;
            }
            // check if there is id and username
            if(data.hasOwnProperty('id') && data.hasOwnProperty('username')) {
                $scope.friend = data;
                $scope.foundFriend = true;
                // DONT ADD YET
                // $scope.friends.push(data);
            }
        }, function(response) {
            // ERROR CALLBACK
            $scope.error = "Cannot access the server at this time...";
        });
    }

    $scope.addFriend = () => {
        console.log("addfriend()");
        if ($scope.friend == null) {
            $scope.error = "No friend selected";
            return;
        }
        // show the loader...
        $scope.addLoader = true;
        $http.post("/api/add/friend/" + $scope.friend.id, {}).then(function(response){
            // SUCCESS CALLBACK...
            console.log("friend:", $scope.friend);
            data = response.data;
            if (data.hasOwnProperty('status') && data.status == 'error') {
                $scope.error = data.error;
                $scope.addLoader = false;
            } else if (data.hasOwnProperty('status') && data.status == 'ok') {
                // message is not empty...
                $scope.chat[$scope.friend.username] = []
                $scope.friends.push($scope.friend);
                $scope.friend = null;
                $scope.friendName = "";
                $scope.addLoader = false;
            }
        }, function(response) {
            // ERROR CALLBACK...
            $scope.error = "Cannot access the server at this time...";
        });
    }

    $scope.logout = () => {
        console.log("logout()");
        $cookies.remove("deewebchat", {path: "/"});
        $window.location.href = "/";
        // TODO :: Send the server a message of logout so that the user can be removed from the DB
        // $http.post("/api/logout", {}).then(function(response){
        //     // successful request...

        // }, function(response) {
        //     // error request...
        // });
    }

});