userarea = angular.module("userarea", []);

userarea.controller("user", function($scope, $location, $http) {
    urlsplit = $location.absUrl().split("/");
    username = urlsplit[urlsplit.length - 1];
    $scope.username = username;
    $scope.friends = [{id: 60, username: "Silvia"}];
    // if a user search returns a result...
    $scope.foundFriend = false;
    $scope.friendName = "";

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
        $scope.friends.push($scope.friend);
        $scope.friend = null;
        $scope.friendName = "";
    }
});