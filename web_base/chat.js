$(document).ready(function() {
	chatwin = $("#chatwindow")
	chatwin.css("bottom",$("#inputbox").height()+60)
	$("#inputbox").keyup(function(evt) {
		chatwin.animate({"bottom": $("#inputbox").height()+60},100, "swing", function() {
			chatwin.children("div").scrollTo("max", 100);
		});
	});
});

//if a token is not set redir to login
if(Cookies.get("token")===undefined) {
  window.location = "index.html";
}
var url="://websocket"
if (location.protocol === 'https:') {
    url = "wss" + url
} else {
	url = "ws" + "url"
}
var ws=null
function connect() {
	ws = new WebSocket(url)
}();
ws.onclose = function() {
	//try to reconnect
	connect();
}
