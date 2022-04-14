var qs = Qs;

const image_select = document.getElementById("image_select");
const access = document.getElementById("access");
const pinned = document.getElementById("pinned");
const tags = document.getElementById("tags");
const image_cont = document.getElementById("image_cont");
const err = document.getElementById("err");
const loading_info = document.getElementById("loading_info");
const upload_btn = document.getElementById("upload_btn");
const loader_init = document.getElementById("loader_init");

axios.post("/fetch_user").then(res => {
    switch (res.data.IsVerified) {
        case "Yes":
            break;
        case "No":
            window.location.assign("/verify");
            break;
        default:
            break;
    }
});

function Upload() {
    window.scrollTo({
        top: 0,
        left: 0,
        behavior: 'smooth'
    });
    if (image_select.files.length > 0) {
        upload_btn.setAttribute('disabled', 'disabled');
        loading_info.innerText = "";
        loader_init.style.display = "block";
        const refreshRate = 1000 / 50;
        const maxXPosition = 85;
        let speedX = 1;
        let positionX = 0;

        window.setInterval(() => {
            positionX = positionX + speedX;
            if (positionX > maxXPosition || positionX < 0) {
                speedX = speedX * (-1);
            }
            loader_init.style.left = positionX + '%';
        }, refreshRate);

        let uploaded = [];
        for (var i = 0; i < image_select.files.length; i++) {
            function Config(filename) {
                axios.post('/update_config', qs.stringify({
                    Access: access.value,
                    Pinned: pinned.value.trim(),
                    Tags: tags.value.trim().toLowerCase(),
                    Filename: filename
                }), {
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded'
                    }
                }).then(res => {
                    uploaded.push(res.data.Name);
                    if (uploaded.length == image_select.files.length) {
                        upload_btn.setAttribute('disabled', 'disabled');
                        loading_info.innerText = "Upload complete"
                        loader_init.style.display = "none";
                        setTimeout(() => {
                            window.location.reload();
                        }, 2000);
                    } else {
                        loading_info.innerText = "Uploading... Please be patient"
                    }
                });
            }

            const formData = new FormData();
            formData.append("meemz_upload", image_select.files[i]);

            axios.post('/upload_meemz', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            }).then(res => {
                if (res.data.Name !== "") {
                    Config(res.data.Name);
                }
            }).catch(err => {
                if (err) {
                    throw err;
                }
            });
        }
    } else {
        err.innerHTML = `<small style="background-color : transparent ; font-family : 'Maven Pro', sans-serif ; color : red;">Choose an image to proceed</small>`;
        err.style.animation = "headShake";
        err.style.animationDuration = "800ms";
    }
}

function DisplayImages() {
    $(image_cont).empty();
    for (let i of image_select.files) {
        let reader = new FileReader();
        let img = document.createElement("img");

        img.style.flex = "33%";
        img.style.width = "200px";
        img.style.height = "200px";
        img.style.margin = "10px";
        img.style.borderRadius = "6px";

        reader.onload = () => {
            img.setAttribute("src", reader.result);
        }
        image_cont.appendChild(img);
        reader.readAsDataURL(i);
    }
}