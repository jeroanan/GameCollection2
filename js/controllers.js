var gameCollectionControllers = angular.module('gameCollectionControllers', []);

gameCollectionControllers.controller('AddGameController', function($scope, $http) {

  $http.get('/json/getplatforms').then(function (data) {
    $scope.platforms = data.data.Platforms;
  });

  $http.get('/json/getgenres').then(function (data) {
    $scope.genres = data.data.Genres;
  });

  $scope.game = {
    'title': '',
    'platform': '',
    'numberowned': 0,
    'numberboxed': 0,
    'numberofmanuals': 0,
    'datepurchased': '',
    'approximatepurchasedate': false,
    'genre': '',
    'notes': ''
  };

  $scope.saveGame = function() {
    $http.post('/json/savegame', JSON.stringify($scope.game)).then(function successCallback(response) {

    }, function errorCallback(response) {
      console.log("Error while saving game.")
    });
  }
});

gameCollectionControllers.controller("AllGamesController", function($scope, $http) {

  $http.get('/json/getgames').then(function(data) {
    $scope.games = data.data.Games;
    $scope.numGames = data.data.Games.length;
  });

  $scope.deleteGame = function(game) {
    $http.post('/json/deletegame', game);
  }
});
