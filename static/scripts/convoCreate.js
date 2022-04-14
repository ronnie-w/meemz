var qs = Qs;

const path = window.location.pathname;
const locale = path.slice(14, path.length);
const link = locale.split("%20").join(" ");

const image_cont = document.getElementById("image_cont");
const convo_topic = document.getElementById("convo_topic");
const convo_title = document.getElementById("convo_title");
const convo_maxnum = document.getElementById("convo_maxnum");
const image_select = document.getElementsByClassName("input_image_select")[0];
const err = document.getElementById("err");

let auth_key = document.cookie.match(new RegExp('(^| )uid=([^;]+)'))[2];

if (link !== "new") {
    convo_topic.value = link;
    convo_topic.setAttribute("readonly", "readonly");
}

function error(msg) {
    err.style.display = "block";
    err.innerHTML = `<small style=backgroundColor : "transparent" , fontFamily : "'Maven Pro', sans-serif" , color : "red"}}>${msg}</small>`;
    err.style.color = "red";
    err.style.animation = "headShake";
    err.style.animationDuration = "800ms";
}

function DisplayBanner() {
    $(image_cont).empty();
    let img = document.createElement('img');

    //image styling
    img.style.width = "200px";
    img.style.height = "200px";
    img.style.borderRadius = "10px";

    let reader = new FileReader();

    reader.onload = () => {
        img.setAttribute("src", reader.result);
    }

    $(image_cont).append(img);
    switch (image_select.files.length) {
        case 0:
            $(image_cont).empty();
            break;

        default:
            reader.readAsDataURL(image_select.files[0]);
            break;
    }
}

function CreateConvo() {
    if (convo_topic.value.length === 0 || convo_title.value.length === 0 || convo_maxnum.value.length === 0) {
        error("Invalid info! Check your input");
    } else if (image_select.files.length === 1) {
        const formData = new FormData();
        formData.append("convo-banner", image_select.files[0]);

        axios.post("/convo_banner_upload", formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        }).then(res => {
            PushToDb(res.data.Msg);
        }).catch(err => {
            if (err) throw err
        });

        function PushToDb(filename) {
            axios.post("/create_convo", qs.stringify({
                Topic: convo_topic.value.trim().toLowerCase(),
                Title: convo_title.value.trim(),
                MaxMembers: convo_maxnum.value.match(/[0-9]/gi).join(""),
                FileName: filename
            }), {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            }).then((res) => {
                window.location.assign(`https://meemzchat.cf/main/details=${auth_key},${res.data.Value}`)
            });
        }
    } else if (image_select.files.length === 0) {
        error("Choose a converstion banner to proceed");
    }
}