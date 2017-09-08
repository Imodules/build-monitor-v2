'use strict';

require('font-awesome/css/font-awesome.css');
require('./styles/main.scss');

const Elm = require('./Main.elm');
const mountNode = document.getElementById('main');

const app = Elm.Main.embed(mountNode, window.options);

const TOKEN_KEY = "authToken";

app.ports.setTokenStorage.subscribe(function(token) {
	localStorage.setItem(TOKEN_KEY, token);
});

app.ports.getTokenFromStorage.subscribe(function() {
	const token = localStorage.getItem(TOKEN_KEY);
	if (token) {
		app.ports.gotTokenFromStorage.send(token);
	}
});

app.ports.logout.subscribe(function() {
	localStorage.removeItem(TOKEN_KEY);
	window.location = "/";
});