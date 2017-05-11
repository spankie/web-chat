var app = angular.module("webchat", []);

app.controller("ctrl1", function($scope){
    $scope.username = "";
    $scope.password = "";
    // $scope.message = "";
    $scope.check = function() {
        if ($scope.username != "" && $scope.password != "") {
            $scope.message = $scope.username + " : " + $scope.password;
        } else {
            $scope.message = "error"
        }
    }
});