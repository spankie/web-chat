var app = angular.module("webchat", []);

app.controller("login", function($scope){
    $scope.username = "";
    $scope.password = "";
    // $scope.message = "";
    $scope.login = function() {
        if ($scope.username != "" && $scope.password != "") {
            $scope.error = "";
            // $scope.message = $scope.username + " : " + $scope.password;
        } else {
            $scope.error = "Username and password cannot be empty";
        }
    }
});