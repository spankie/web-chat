var app = angular.module("webchat", ['ngCookies']);

app.controller("login", function($scope, $http, $cookies, $window){
    $scope.username = "";
    $scope.password = "";
    $scope.message = "";
    $scope.login = function() {

        if ($scope.username != "" && $scope.password != "") {
            $http.post("/api/login", {username: $scope.username, password: $scope.password}).then(function(response){
                // SUCCESS CALLBACK:
                data = response.data;
                console.log("response: ", data);
                if(data.status == "ok") {
                    console.log("You are logged in");
                    $scope.error = "";
                    // set the cookie here
                    $cookies.put("deewebchat", data.cookie);
                    // navigate to the chat page...
                    // $location.url("http://localhost:8080/" + $scope.username);
                    $window.location.href = "/web/" + $scope.username;
                } else {
                    // display error here
                    console.log("You are NOT logged in.");
                    $scope.error = "Username or password is incorrect.";
                }

            }, function(response) {
                // ERROR CALLBACK:
                console.log("error response: ", response);                
                $scope.error = "Error response";
            })

        } else {
            $scope.error = "Username and password cannot be empty";
        }
    }
});