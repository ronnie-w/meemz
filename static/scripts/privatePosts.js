var qs = Qs;

const private_meemz = document.getElementById("private_meemz");

let content;

axios.post("/fetch_private_posts").then(res => {
    if (res.data === null) {
        private_meemz.innerHTML = "<p style='top : 20%; position : absolute; display : grid; place-content : center; width : 100%'>You have no private Meemz</p>";
    } else {
        content = res.data;
        content.forEach(c => {
            $("#private_meemz").prepend(
                `<div class="private_meem_main_div ${c.ImgName}">
                <div class="private_meem">
                    <img src="/static/uploads/${c.ImgName}" alt="private_meem" class="image"
                        style="border-radius : 6px" />
                    <div class="private_actions">
                        <button class="private_to_public" onclick="Publicize('${c.ImgName}')"><i class="far fa-globe private_globe"></i></button>
                        <button class="private_delete" onclick="Delete('${c.ImgName}')"><i class="far fa-trash-alt private_trash"></i></button>
                    </div>
                </div>
            </div>`
            );
        });
    }
});

function Publicize(imgName) {
    let private_meem_div = document.getElementsByClassName(`private_meem_main_div ${imgName}`)[0];
    axios.post("/publicize_post", qs.stringify({
        ImgName: imgName
    }), {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }).then(_res => {
        private_meem_div.style.animation = "fadeOut";
        private_meem_div.style.animationDuration = "800ms";
        setTimeout(() => {
            private_meem_div.style.display = "none";
        }, 799);
    });
}

function Delete(imgName) {
    let private_meem_div = document.getElementsByClassName(`private_meem_main_div ${imgName}`)[0];
    axios.post("/delete_post", qs.stringify({
        ImgName: imgName
    }), {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }).then(_res => {
        private_meem_div.style.animation = "fadeOut";
        private_meem_div.style.animationDuration = "800ms";
        setTimeout(() => {
            private_meem_div.style.display = "none";
        }, 799);
    });
}