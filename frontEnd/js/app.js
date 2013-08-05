(function() {

  var App = Ember.Application.create();
  
  App.Router.map(function() {
    this.route('articles');
	this.route('photos');
	this.route('login');
  });
  
  App.AuthenticatedRoute = Ember.Route.extend({
	  getJSONWithToken: function(url){
  		var token = this.controllerFor('login').get('token');
  		return $.getJSON(url, { token: token })
	  },
	  events: {
		  error: function(reason, transition){
			  if(reason.status == 401){
				  alert('You must log in!');
				  var loginController = this.controllerFor('login');
				  loginController.set('attemptedTransition', transition);
				  this.transitionTo('login');
			  } else {
				  alert('Something went wrong!');
			  }
		  }
	  }
  });
  
  App.ArticlesRoute = App.AuthenticatedRoute.extend({
    model: function() {
		return this.getJSONWithToken('/api/articles.json');
    }
  });
  
  App.PhotosRoute = App.AuthenticatedRoute.extend({
    model: function() {
		return this.getJSONWithToken('/api/photos.json');
    }
  });
  
  App.LoginRoute = Ember.Route.extend({
    setupController: function(controller, context) {
      controller.reset();
    }
  });
  
  App.LoginController = Ember.Controller.extend({
	  
	  reset: function() {
	    this.setProperties({
	      username: "",
	      password: "",
	      errorMessage: ""
	    });
	  },
	  
	  login: function() {
		  
		  var self = this, data = this.getProperties('username', 'password');
		  
	      // Clear out any error messages.
	      this.set('errorMessage', null);
		  
		  $.ajax({
		    url: '/api/auth.json',
		    type: 'POST',
		    contentType: 'application/json; charset=utf-8',
		    data: JSON.stringify(data)
		    // you can add callbacks as “complete”, “success”, etc.
		  }).then(function(response) {
			  // Check the response for the token
			  self.set('errorMessage', response.message);
			  
			  if(response.success){
				  alert('Login Succeeded!');
				  self.set('token', response.token);
				  var attemptedTransition = self.get('attemptedTransition');
				  if(attemptedTransition){
					  attemptedTransition.retry();
					  self.set(attemptedTransition, null);
				  } else {
					  self.transitionToRoute('articles');
				  }
				  self.transitionTo('articles');
			  }
		  });
	  }
  });
  

})();
