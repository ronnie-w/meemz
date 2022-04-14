var qs = Qs;

const profile_top_bar = document.getElementById("profile_top_bar");
const profile_config_main = document.getElementById("profile_config_main");
const posts = document.getElementById("posts");
const my_upload = document.getElementById("my_upload");
const searched_rooms = document.getElementById("searched_rooms");
const bio = document.getElementById("bio");

let user_data;
let subs;
let uploads;
let rooms;

let auth_key = document.cookie.match(new RegExp('(^| )uid=([^;]+)'))[2];

axios.post("/fetch_user").then(res =>{
    user_data = res.data;

    if (user_data.Bio === "No bio found") {
        bio.style.display = "none";
    }
    
    $("#profile_img").attr("src", `/static/profile-pictures/${user_data.ProfileImg}`);
    $(".profile_username").text(user_data.Username);
    $(".users_bio").text(user_data.Bio);

    axios.post("/fetch_profile_rooms" , qs.stringify({
        Username : user_data.Username
    }) , {
        headers : {
            'Content-Type' : 'application/x-www-form-urlencoded'
        }
    }).then(res =>{
        if (res.data !== null) {
            rooms = res.data;
        }

        rooms.forEach(room =>{
            $("#searched_rooms").append(
                `<div style="border : 1px solid rgb(48, 47, 47); margin-top : 10px; box-shadow: .2px .2px 5px grey;">
                    <center>
                        <img src="/static/convo_banners/${room.TopicProfile}" alt="convo_banner"
                            style="width : 100px; height : 100px; border-radius : 10px; margin : 5px; box-shadow : .2px .2px 3px wheat;" />
                    </center>
                    <center>
                        <p style="margin-bottom : 10px">${room.Title}</p>
                        <button style="color : wheat; background-color : grey;" onclick="javascript:window.location.assign('https://meemzchat.cf/main/details=${auth_key},${room.ChatRoomId}')">Join</button>
                    </center>
                </div>`
            );
        });
    });

    axios.post("/fetch_subs" , qs.stringify({
        Username : user_data.Username
    }) , {
        headers : {
            'Content-Type' : 'application/x-www-form-urlencoded'
        }
    }).then(res =>{
        if (res.data.Subscribers === 1) {
            subs = `${res.data.Subscribers} subscriber`;
        }else{
            subs = `${res.data.Subscribers} subscribers`;
        }

        $("#profile_subs").text(subs);
    });
});

axios.post("/my_uploads").then(res =>{
    if (res.data !== null) {
        uploads = res.data.reverse();
    }

    uploads.forEach(upload =>{
        $("#my_upload").append(
            `<div>
                <a href="/image_stats/${upload.ImgName}"><img class="upload" src="/static/uploads/${upload.ImgName}" alt="my_upload" style="width : 115px ; height : 115px ; border : 1px solid rgb(48, 47, 47);"loading="lazy" /></a>
            </div>`
        );
    });
});

window.addEventListener("scroll" , () =>{
    if (window.scrollY > 185) {
        profile_top_bar.style.display = "block";
        profile_top_bar.innerHTML = `<a href='/' style='float:left;margin:10px'><i class='fas fa-chevron-left'></i></a><p style='margin:7px'>${user_data.Username}</p>`;
    }else if(window.scrollY < 185){
        profile_top_bar.style.display = "flex";
        profile_top_bar.innerHTML = `<a href='/profile_edit'><i class='profile_edit fal fa-user-edit'></i></a><a href='/private_posts'><i class='private fal fa-key'></i></a><i class="sign_out fal fa-sign-out-alt" onclick="javascript: document.cookie = 'uid=0; expires=Thu, 19 Dec 2002 12:00:00 UTC'; window.location.replace('/signup')"></i>`;
    }
    if (window.scrollY > 275){
        profile_config_main.style.position = "fixed";
        profile_config_main.style.top = "0";
        profile_config_main.style.marginTop = "37px";
    }else{
        profile_config_main.style.position = "relative";
        profile_config_main.style.marginTop = "10px";
    }

    document.cookie = `screen_pos=${window.scrollY};`;
});

function Posts() {
    axios.post("/my_uploads").then(res =>{
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

    axios.post("/fetch_profile_rooms" , qs.stringify({
        Username : user_data.Username
    }) , {
        headers : {
            'Content-Type' : 'application/x-www-form-urlencoded'
        }
    }).then(res =>{
        if (res.data !== null) {
            rooms = res.data;
        }
        if (class_name === "far fa-images") {
            my_upload.style.display = "flex";
            searched_rooms.style.display = "none";
            img_btn.style.color = "wheat";
            rooms_btn.style.color = "white";
    
            Posts();
        } else {
            my_upload.style.display = "none";
            searched_rooms.style.display = "flex";
            rooms_btn.style.color = "wheat";
            img_btn.style.color = "white";
    
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

function Logout() {
    document.cookie = 'uid=0; expires=Thu, 19 Dec 2002 12:00:00 UTC';
    window.location.replace('/signup');
}