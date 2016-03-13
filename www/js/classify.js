$("document").ready(function(){

	$('#menu-list').on('click', function() {
		$('#menu-bottom').toggleClass('custom-menu-tucked');
		$('#menu-list').toggleClass('x');
	});

	var ws = start('ws://localhost:8080/ws');

	$('input#new-directory').on('change', function(){
		send(ws, "new-directory", {
			path: $('input#newDirectory').val()
		});
	});
});

var onEvent = {
	newFile: function(data) {
		$("input#new-directory").after($("<p>").text(data.FullPath));
	},
}

function start(url)
{
	ws = new WebSocket(url);

	ws.onopen = function(e) {
		console.log('Connection OK');
	};

	var onClose = {
		1000: "Close OK",
		1001: "Server/Client DOWN",
		1002: "Protocol error",
		1003: "Type of data not accepted",
		1004: "Reserved",
		1005: "No status code",
		1006: "Connection was closed abnormally",
		1007: "Message not consistent",
		1008: "Message violates policy",
		1009: "Message too big",
		1010: "Negotiation failed",
		1011: "Unexpected condition",
		1015: "TLS handshake failed"
	};

	ws.onclose = function(e) {
		console.log('Connection close: %s', onClose[e.code]);
		setTimeout(function(){start(url)}, 5000);
	}

	ws.onmessage = function(e) {
		var obj = JSON.parse(e.data);
		console.log(' --> %o', obj);
		onEvent[obj.Type](obj.Data);
	};

	ws.onerror = function(e) {
		console.log(' !!! %o', e);
	};

	return ws
}

function send(ws, type, data) {
	var json = JSON.stringify({type: type, data: data});
	console.log(' <-- %s', json);
	ws.send(json);
}
