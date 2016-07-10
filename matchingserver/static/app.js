var ws = new WebSocket("ws://" + document.location.host + "/echo");
ws.onmessage = function(event) {
    console.log(event);
    var msgListBox = document.querySelector("#msgListBox");
    var line = document.createElement("div")
    line.textContent = event.data;
    msgListBox.appendChild(line);
};

ws.onclose = function(event) {
    console.log(event);
    alert("Chat disconnected please refresh");
};

var formMsg = document.querySelector("#formMsg");
formMsg.addEventListener("submit", function(event) {
    var msgBox = document.querySelector("#msgBox");
    var msgListBox = document.querySelector("#msgListBox");
    var text = msgBox.value;
    console.log(text);
    var line = document.createElement("div")
    line.textContent = text;
    msgListBox.appendChild(line);
    ws.send(text + "\n");
    msgBox.value = "";
    event.preventDefault();
    return false;
});