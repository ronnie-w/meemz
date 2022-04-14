var qs = Qs;

const browse_topics = document.getElementById("browse_topics");
const convo_topics = document.getElementById("convo_topics");
const convo_search_rooms = document.getElementById("convo_search_rooms");
const searched_rooms = document.getElementById("searched_rooms");

let auth_key = document.cookie.match(new RegExp('(^| )uid=([^;]+)'))[2];
let rooms;

window.addEventListener("scroll", () => {
    if (window.scrollY > 45) {
        browse_topics.style.display = "block";
    } else {
        browse_topics.style.display = "none";
    }
});

function Loader() {
    if (convo_search_rooms.value.length > 0) {
        $(searched_rooms).empty();
        searched_rooms.style.display = "flex";
        convo_topics.style.display = "none";
        axios.post("/search_room_titles", qs.stringify({
            ChatRoom: convo_search_rooms.value.trim()
        }), {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }).then(res => {
            if (res.data !== null) {
                rooms = res.data;

                rooms.forEach(room =>{
                    $(searched_rooms).append(
                        `<div style="border : 1px solid rgb(48, 47, 47) ; margin-top : 10px; box-shadow: .2px .2px 5px grey;">
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
            }
        });
    } else {
        searched_rooms.style.display = "none";
        convo_topics.style.display = "block";
    }
};