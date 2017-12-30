var socket;
$("#connect").click(function (event) {
    socket = new WebSocket("ws://127.0.0.1:8000");
    socket.onopen = function () {
        alert("Socket has been opened");
    }
    socket.onmessage = function (msg) {
        alert(msg.data);
    }
    socket.onclose = function () {
        alert("Socket has been closed");
    }
});
$("#send").click(function (event) {
    socket.send("send from client");
});
$("#close").click(function (event) {
    socket.close();
})