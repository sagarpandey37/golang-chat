$(function () {
    let websocket = new WebSocket("ws://" + window.location.host + "/websocket?userID=32");
    let room = $("#chat-text");
    websocket.addEventListener("message", function (e) {
        let data = JSON.parse(e.data);
        let chatContent = `<p><strong>${data.username}</strong>: ${data.text}</p>`;
        room.append(chatContent);
        room.scrollTop = room.scrollHeight; // Auto scroll to the bottom
    });
    $("#input-form").on("submit", function (event) {
        event.preventDefault();
        let username = $("#input-username")[0].value;
        let text = $("#input-text")[0].value;
        websocket.send(
            JSON.stringify({
                "channelKey":1,
                "ChannelType":0,
                "sender": {
                            "UserName" : "Sagar Pandey",
                            "UserID" : 32,
                            "userStatus" : true,
                            "userLastActivityTime" : "10 Feb,2021 11:30:00"
                          },
                          
                "reciever" : {
                            "UserName" : "Tom",
                            "UserID" : 16,
                            "userStatus" : true,
                            "userLastActivityTime" : "10 Feb,2021 11:30:00"
                          },
                "CreateDate": "10 Feb,2021 11:30:00",
                "message": {
                    "text" : `${text}`,
                    "author" : {
                        "UserName" : "Sagar Pandey",
                        "UserID" : 32,
                        "userStatus" : true,
                        "userLastActivityTime" : "10 Feb,2021 11:30:00"
                      },
                    "createDate": "10 Feb,2021 11:30:00"
                }
                })
        );
        $("#input-text")[0].value = "";
    });
});