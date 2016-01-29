var gameCollectionServices = angular.module('gameCollectionServices', ['ngResource']);

gameCollectionServices.factory('GetGenres', ['$resource',
  function($resource) {
    return $resource('/json/getgenres', {}, {
      query: {mehtod: 'GET', params: {} }
    });
  }]);

gameCollectionServices.factory('GetPlatforms', ['$resource',
  function($resource) {
    return $resource('/json/getplatforms', {}, {
      query: {mehtod: 'GET', params: {} }
    });
  }]);
