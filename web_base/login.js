$(document).ready(function() {
  //hide NO-JS prompt
  if(Cookies.get("status")===undefined) {
    $("#statusbox").hide();
  } else {
    $("#statusbox h2").text(decodeURIComponent(Cookies.get("status").replace(/\+/g,  " ")));
    Cookies.remove("status",  { path: '' });
  }
});


//if a token is already set, login shoud not be needed
if(Cookies.get("token")!==undefined) {
  window.location = "cookieredir.php";
}
