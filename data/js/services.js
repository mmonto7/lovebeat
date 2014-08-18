var lovebeatServices = angular.module('lovebeatServices', ['ngResource']);

lovebeatServices.factory('Service', ['$resource',
  function($resource){
    return $resource('/api/services/?view=:viewId', {}, {
      query: {method:'GET', isArray:true}
    });
  }]);

lovebeatServices.factory('View', ['$resource',
  function($resource){
    return $resource('/api/views/?', {}, {
      query: {method:'GET', isArray:true}
    });
  }]);
