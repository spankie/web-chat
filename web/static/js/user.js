userarea = angular.module("userarea", []);

userarea.controller("user", function($scope) {
    $scope.username = "spankie";
    $scope.foundFriend = "false";
});