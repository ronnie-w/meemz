var qs = Qs;

let notifications;

axios.post("/notifications_go").then(res => {
    if (res.data === null) {
        notifications_div.innerHTML = "<div style='width : 100%; display : grid; place-content : center; margin-top : 50px;'><p>You have no new notifications</p></div>"
    } else {
        notifications = res.data;

        notifications.forEach(n =>{
            $("#notifications_div").append(
                `<div style="margin-bottom : 50px ; margin-top : 10px;" class="n_div">
                <h6 style="margin-left : 5px;">Received on ${n.ReceiveTime}</h6>
                <img src="/static/profile-pictures/${n.ProfileImg}" alt="notification_banner" style="width : 35px ; height : 35px ; border-radius : 200px ; float : left ; margin-left : 5px ; margin-top : 7px ; border : 2px solid white;" />
                <p style="color : grey ; float : left ; background-color : transparent ; font-size : 12px ; margin-top : 7px ; margin-left : 5px ; text-shadow : .5px .5px 7px grey;">${n.Username}</p>
                <br/>
                <div style="width : 100% ; background-color : transparent;">
                    <p style="background-color : transparent ; color : white ; font-size : 15px ; max-width : 300px ; margin-top : 5px ; word-wrap : break-word ; margin-left : 48px;">${n.Notify}</p>
                </div>
                </div>`
            );
        });
    }
});

window.addEventListener("DOMContentLoaded", () =>{
    axios.post("/delete_notifications");
});