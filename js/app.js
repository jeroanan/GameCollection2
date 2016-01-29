var gameCollection = angular.module('gameCollection', [
  'gameCollectionControllers',
  'ngRoute',
  'gameCollectionServices'
]);

gameCollection.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/addgame', {
        templateUrl: '/view/addgame'
      }).
      when('/', {
        templateUrl: '/view/allgames'
      }).
      when('/editgame/:gameid', {
        templateUrl: '/view/editgame'
      })
    }
  ]);
