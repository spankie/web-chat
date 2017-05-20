userarea = angular.module("userarea", []);

userarea.controller("user", function($scope, $location) {
    urlsplit = $location.absUrl().split("/");
    username = urlsplit[urlsplit.length - 1];
    $scope.username = username;
    // if a user search returns a result...
    $scope.foundFriend = true;
    $scope.friend = "";

    $scope.searchFriend = function() {
        // search for a friend from the server.
        // /api/search/friend
    }
});