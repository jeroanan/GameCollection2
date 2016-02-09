var gameCollectionControllers = angular.module('gameCollectionControllers', []);

gameCollectionControllers.controller('AddGameController', ['$scope', '$http', 'GetPlatforms', 'GetGenres',

  function($scope, $http, GetPlatforms, GetGenres) {

    $scope.platforms = GetPlatforms.query();
    $scope.genres = GetGenres.query();

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
      $scope.hasError = false;
      $scope.errorMessage = '';

      $http.post('/json/savegame', JSON.stringify($scope.game)).then(function successCallback(response) {

      }, function errorCallback(response) {
        $scope.hasError = true;
        $scope.errorMessage = 'Error while saving game: ' + response.data;
      });
    }

    $scope.hasError = false;
    $scope.errorMessage = '';
}]);

gameCollectionControllers.controller("AllGamesController", function($scope, $http) {

  $http.get('/json/getgames').then(function successCallback(data) {
    $scope.games = data.data.Games;
    $scope.numGames = data.data.Games.length;
  }, function errorCallback(response) {
    $scope.hasError = true;
    $scope.errorMessage = 'Error while retrieving games: ' + response.data;
  });

  $scope.deleteGame = function(game) {
    $http.post('/json/deletegame', game);
  }

  $scope.sortField = 'Title';
  $scope.reverseSort = false;
  $scope.hasError = false;
  $scope.errorMessage = '';

});

gameCollectionControllers.controller("EditGameController", ['$scope', '$http', '$routeParams', 'GetPlatforms',
  'GetGenres',
  function($scope, $http, $routeParams, GetPlatforms, GetGenres) {
    $scope.platforms = GetPlatforms.query();
    $scope.genres = GetGenres.query();

    $http.get('/json/getgame/' + $routeParams.gameid).then(function(data) {
      $scope.game = data.data;
    });

    $scope.saveGame = function() {
      var scopeGame = $scope.game;
      var game = {
        'title': scopeGame.Title,
        'platform': scopeGame.PlatformId,
        'numberowned': scopeGame.NumberOwned,
        'numberboxed': scopeGame.NumberBoxed,
        'numberofmanuals': scopeGame.NumberOfManuals,
        'datepurchased': scopeGame.DatePurchased,
        'approximatepurchasedate': scopeGame.ApproximatePurchaseDate,
        'genre': scopeGame.GenreId,
        'notes': scopeGame.Notes,
        'RowId': scopeGame.RowId
      };

      $http.post('/json/savegame', JSON.stringify(game)).then(function successCallback(response) {
        $scope.hasError = false;
        $scope.errorMessage = '';
      }, function errorCallback(response) {
        $scope.hasError = true;
        $scope.errorMessage = 'Error while saving game: ' + response.data;
      });
    }

    $scope.cancel = function() {
      window.location = '/';
    }

    $scope.hasError = false;
    $scope.errorMessage = '';
}]);

gameCollectionControllers.directive('errorbox', function() {
  return {
    restrict: 'E',    
    templateUrl: '/directives/error.html'
  }
});
