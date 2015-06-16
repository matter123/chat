$(document).ready(function() {
	chatwin = $("#chatwindow")
	chatwin.css("bottom",$("#inputbox").height()+60)
	$("#inputbox").keyup(function(evt) {
		if (evt.which !=13) {
			chatwin.animate({"bottom": $("#inputbox").height()+60},100, "swing", function() {
				chatwin.children("div").scrollTo("max", 100);
			});
		} else {
			ws.send(JSON.stringify({time: Math.floor(Date.now()/1000).toString(),message: {user: 'test', message: $(this).text()}}))
			$(this).text('')
		}
	});
});

//if a token is not set redir to login
if(Cookies.get("token")===undefined) {
  //window.location = "index.html";
}
var url="://" + window.location.host + "/websocketmock"
if (location.protocol === 'https:') {
    url = "wss" + url
} else {
	url = "ws" + url
}
var ws=null
function connect() {
	console.log(url)
	ws = new WebSocket(url);
}
connect();
ws.onmessage = function(messageevt) {
	message = JSON.parse(messageevt.data)
	console.log(message)
	if ('ping' in message) {
		console.log('ping');
		console.log(JSON.stringify({time: Math.floor(Date.now()/1000).toString(), pong: {nanotime: message.ping.nanotime}}));
		this.send(JSON.stringify({time: Math.floor(Date.now()/1000).toString(), pong: {nanotime: message.ping.nanotime}}));
	}
	if ('message' in message) {
		console.log('message');
		$('#chats').append('<div class="chat"><span class="name">' + message.message.user + ': </span>' + message.message.message + '</div>')
	}
	if('leave' in message) {
		$('#chats').append('<div class="chat"><span class="name">' + message.leave.user + ': </span>Has left.</div>')
	}
	if('join' in message) {
		$('#chats').append('<div class="chat"><span class="name">' + message.join.user + ': </span>Has Joined.</div>')
	}
}
