userarea = angular.module("userarea", ['ngCookies']);

userarea.controller("user", function($scope, $location, $cookies, $http, $window) {
    var chatBox = document.getElementById("chatbox");
    chatBox.scrollTop = chatBox.scrollHeight;
    
    urlsplit = $location.absUrl().split("/");
    username = urlsplit[urlsplit.length - 1];
    $scope.username = username;
    $scope.friends = [];
    $http.post("/api/get/friends", {}).then(function(response){
        // success response callback
        data = response.data;
        if (data.hasOwnProperty('status') && data.status == 'error') {
            // You have no friends
            // tell the user
            console.log("no friends");
        } else if (data.hasOwnProperty('status') && data.status == 'ok') {
            $scope.friends = data;
            console.log("Some friends");
        }
        console.log("scope test");
    }, function(response) {
        // error response callback
    });

    // $scope.friends = [{id: 60, username: "Silvia"}];
    // if a user search returns a result...
    $scope.foundFriend = false;
    $scope.friendName = "";
    $scope.friend = null;
    $scope.addLoader = false;

    $scope.searchFriend = function() {
        console.log("FRIEND: ", $scope.friend)
        // search for a friend from the server.
        // /api/search/friend
        $http.post("/api/search/friend", {username: $scope.friendName}).then(function(response){
            // SUCCESS CALLBACK
            data = response.data;
            console.log(data);
            if (data.hasOwnProperty('status') && data.status == 'error') {
                // Could not get the user...
                // Display the error here
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
                $scope.friends.push($scope.friend);
                $scope.friend = null;
                $scope.friendName = "";
                $scope.addLoader = false;
            }
        }, function(response) {
            // ERROR CALLBACK...
            $scope.error = "Cannot access the server at this time..."
        });
    }

    $scope.logout = () => {
        console.log("logout()");
        $cookies.remove("deewebchat", {path: "/"});
        $window.location.href = "/";
        // $http.post("/api/logout", {}).then(function(response){
        //     // successful request...

        // }, function(response) {
        //     // error request...
        // });
    }

});