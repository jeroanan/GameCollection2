var gameCollection = angular.module('gameCollection', [
  'gameCollectionControllers',
  'ngRoute'
]);

gameCollection.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/addgame', {
        templateUrl: '/view/addgame'
      })
    }
  ]);
