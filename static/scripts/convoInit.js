var qs = Qs;

const searched_rooms = document.getElementById("searched_rooms");
const no_rooms = document.getElementById("no_rooms");
const num_indicator = document.getElementById("num_indicator");

const path = window.location.pathname;
const locale = path.slice(12, path.length);
const link = locale.split("%20").join(" ");

let auth_key = document.cookie.match(new RegExp('(^| )uid=([^;]+)'))[2];
let rooms;

$("#start_a_convo").attr("href", `/convo_create/${link}`);

axios.post("/search_room_topics", qs.stringify({
    ChatRoom: link.trim().toLowerCase()
}), {
    headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
    }
}).then(res => {
    if (res.data !== null) {
        rooms = res.data
        rooms.forEach(room => {
            $(searched_rooms).append(
                `<div style="border : 1px solid rgb(48, 47, 47) ; margin-top : 10px ; box-shadow: .2px .2px 5px grey;">
                <center>
                    <img src="/static/convo_banners/${room.TopicProfile}" alt="convo_banner"
                        style="width : 100px ; height : 100px ; border-radius : 10px ; margin : 5px ; border : 1px solid rgb(48, 47, 47);" />
                </center>
                <center>
                    <p style="margin-bottom : 10px;">${room.Title}</p>
                    <button style="color : wheat ; background-color : grey;" onclick="javascript:window.location.assign('https://meemzchat.cf/main/details=${auth_key},${room.ChatRoomId}')">Join</button>
                </center>
            </div>`
            );
        });
        if (res.data.length === 1) {
            num_indicator.innerHTML = `<p>${res.data.length} room found</p>`;
        } else {
            num_indicator.innerHTML = `<p>${res.data.length} rooms found</p>`;
        }
    } else {
        no_rooms.style.display = "block";
        num_indicator.style.display = "none";
    }
});