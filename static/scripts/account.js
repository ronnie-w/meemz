var qs = Qs;

const profile_top_bar = document.getElementById("profile_top_bar"),
    profile_config_main = document.getElementById("profile_config_main"),
    posts = document.getElementById("posts"),
    my_upload = document.getElementById("my_upload"),
    searched_rooms = document.getElementById("searched_rooms"),
    bio = document.getElementById("bio");

let acc,
    user_data,
    subs,
    uploads,
    rooms,
    numsubs;

const path = window.location.pathname;
const locale = path.slice(6, path.length);
const link = locale.split("%20").join(" ");

let auth_key = document.cookie.match(new RegExp('(^| )uid=([^;]+)'))[2];

$("#my_upload").css({
    marginTop: "5px",
    display: "flex",
    flexWrap: "wrap",
    justifyContent: "space-evenly"
});
$("#searched_rooms").css({
    marginTop: "5px",
    display: "none",
    flexWrap: "wrap",
    justifyContent: "space-evenly",
    marginBottom: "50px"
});

axios.post('/search_users', qs.stringify({
    Username: link
}), {
    headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
    }
}).then((res) => {
    acc = res.data[0];
    $(".profile_subs").append(`<button id="sub_${acc.Username}" style="background-color: ${acc.BGC}; color: ${acc.Color}; font-size: 8px; box-shadow: .5px .5px 7px rgb(255, 166, 0);" onclick="Subscription('sub_${acc.Username}' , '${acc.Username}')">${acc.Subscription}</button>`);
    $("#profile_img").attr("src", `/static/profile-pictures/${acc.ProfileImg}`);
    $(".profile_username").text(acc.Username);
    $(".users_bio").text(acc.Bio);
    if (acc.Bio === 'No bio found') {
        bio.style.display = "none";
    }
}).catch((err) => {
    console.log(err)
});

axios.post("/fetch_subs", qs.stringify({
    Username: link
}), {
    headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
    }
}).then(res => {
    if (res.data.Subscribers === 1) {
        subs = `${res.data.Subscribers} subscriber`;
    } else {
        subs = `${res.data.Subscribers} subscribers`;
    }

    $("#profile_subs").text(subs);
});

axios.post("/fetch_profile_rooms", qs.stringify({
    Username: link
}), {
    headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
    }
}).then(res => {
    if (res.data !== null) {
        rooms = res.data;
    }

    rooms.forEach(room => {
        $("#searched_rooms").append(
            `<div style="border : 1px solid rgb(48, 47, 47); margin-top : 10px; box-shadow: .2px .2px 5px grey;">
                <center>
                    <img src="/static/convo_banners/${room.TopicProfile}" alt="convo_banner"
                        style="width : 100px; height : 100px; border-radius : 10px; margin : 5px; box-shadow : .2px .2px 3px grey;" />
                </center>
                <center>
                    <p style="margin-bottom : 10px">${room.Title}</p>
                    <button style="color : rgb(255, 166, 0); background-color : grey;" onclick="javascript:window.location.assign('https://meemzchat.cf/main/details=${auth_key},${room.ChatRoomId}')">Join</button>
                </center>
            </div>`
        );
    });
});

axios.post("/fetch_user_uploads", qs.stringify({
    Username: link
}), {
    headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
    }
}).then(res => {
    if (res.data !== null) {
        uploads = res.data.reverse();
    }

    uploads.forEach(upload => {
        $("#my_upload").append(
            `<div>
                <a href="/public_stats/${upload.ImgName}"><img class="upload" src="/static/uploads/${upload.ImgName}" alt="my_upload" style="width : 115px ; height : 115px ; border : 1px solid rgb(48, 47, 47);"loading="lazy" /></a>
            </div>`
        );
    });
});

window.addEventListener("scroll", () => {
    if (window.scrollY > 185) {
        profile_top_bar.style.display = "block";
        profile_top_bar.innerHTML = `<a href='/' style='float:left;margin:10px'><i class='fas fa-chevron-left'></i></a><p style='margin:7px'>${link}</p>`;
    } else if (window.scrollY < 185) {
        profile_top_bar.style.display = "flex";
        profile_top_bar.innerHTML = ``;
    }
    if (window.scrollY > 275) {
        profile_config_main.style.position = "fixed";
        profile_config_main.style.top = "0";
        profile_config_main.style.marginTop = "37px";
    } else {
        profile_config_main.style.position = "relative";
        profile_config_main.style.marginTop = "10px";
    }
});

function Posts() {
    axios.post("/fetch_user_uploads", qs.stringify({
        Username: link
    }), {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }).then(res => {
        if (res.data !== null) {
            uploads = res.data.reverse();
        }
        if (uploads.length === 1 && uploads.length !== 0) {
            posts.innerText = `${uploads.length} post found`;
            posts.style.display = "block";
        } else if (uploads.length > 1 && uploads.length !== 0) {
            posts.innerText = `${uploads.length} posts found`;
            posts.style.display = "block";
        }
    });
};

Posts();

function ToggleConfig(class_name) {
    let img_btn = document.getElementsByClassName("posts far fa-images")[0];
    let rooms_btn = document.getElementsByClassName("posts far fa-comment-alt")[0];

    axios.post("/fetch_profile_rooms", qs.stringify({
        Username: link
    }), {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }).then(res => {
        if (res.data !== null) {
            rooms = res.data;
        }
        if (class_name === "far fa-images") {
            my_upload.style.display = "flex";
            searched_rooms.style.display = "none";
            img_btn.style.color = "rgb(255, 166, 0)";
            rooms_btn.style.color = "#121212";

            Posts();
        } else {
            my_upload.style.display = "none";
            searched_rooms.style.display = "flex";
            rooms_btn.style.color = "rgb(255, 166, 0)";
            img_btn.style.color = "#121212";

            if (rooms.length === 1 && rooms.length !== 0) {
                posts.innerText = `${rooms.length} rooms found`;
                posts.style.display = "block";
            } else if (rooms.length > 1 && rooms.length !== 0) {
                posts.innerText = `${rooms.length} rooms found`;
                posts.style.display = "block";
            }
        }
    });
};

function Subscription(btn_id, username) {
    let btn = document.getElementById(btn_id);
    axios.post("/fetch_subs", qs.stringify({
        Username: link
    }), {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }).then(res => {
        numsubs = res.data.Subscribers;
        if (btn.innerText === "Subscribe") {
            axios.post("/subscribe", qs.stringify({
                Username: username
            }), {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            });
            btn.style.backgroundColor = "#121212";
            btn.style.color = "rgb(255, 166, 0)";
            btn.innerText = "Unsubscribe";

            numsubs++;
            if (numsubs === 1) {
                $("#profile_subs").text(`${numsubs} subscriber`);
            } else {
                $("#profile_subs").text(`${numsubs} subscribers`);
            }
        } else {
            axios.post("/unsubscribe", qs.stringify({
                Username: username
            }), {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            });
            btn.style.backgroundColor = "rgb(236, 235, 235)";
            btn.style.color = "#121212";
            btn.innerText = "Subscribe";

            numsubs--;
            if (numsubs === 1) {
                $("#profile_subs").text(`${numsubs} subscriber`);
            } else {
                $("#profile_subs").text(`${numsubs} subscribers`);
            }
        }
    });
};