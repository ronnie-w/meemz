var qs = Qs;

const profile_img_select = document.getElementsByClassName("profile_img_select")[0];
const profile_edit_img = document.getElementById("profile_edit_img");
const profile_edit_username = document.getElementById("profile_edit_username");
const profile_edit_email = document.getElementById("profile_edit_email");
const profile_edit_bio = document.getElementById("profile_edit_bio");
const auth_errs = document.getElementById("auth_errs");
const tip = document.getElementById("tip");

setTimeout(() => {
    $(tip).css({
        animation : "fadeOutUp",
        animationDuration : "800ms"
    });
    setTimeout(() => {
        $(tip).css({
            display : "none"
        });
    }, 800);
}, 5000);

axios.post("/fetch_user").then(res => {
    $("#profile_edit_img").attr("src", `/static/profile-pictures/${res.data.ProfileImg}`);
    $("#profile_edit_username").attr("placeholder", `${res.data.Username}`);
    $("#profile_edit_email").attr("placeholder", `${res.data.Email}`);
    $("#profile_edit_bio").attr("placeholder", `${res.data.Bio}`);
});

function error(msg) {
    auth_errs.style.display = "block";
    auth_errs.innerHTML = `<small style=backgroundColor : "transparent" , fontFamily : "'Maven Pro', sans-serif" , color : "red"}}>${msg}</small>`;
    auth_errs.style.animation = "headShake";
    auth_errs.style.animationDuration = "800ms";
}

function Change() {
    if (profile_img_select.files.length !== 1) return;
    const fileReader = new FileReader();
    fileReader.onload = () => {
        profile_edit_img.setAttribute("src", fileReader.result);
    }
    fileReader.readAsDataURL(profile_img_select.files[0]);

    if (profile_img_select.files.length === 1) {
        let formData = new FormData();
        formData.append("profile_img", profile_img_select.files[0]);
        axios.post("/profile_img_upload", formData, {
            headers: {
                'Content-Type': 'multipart/form-data'
            }
        }).then(res => {
            if (res.data !== null) {
                UpdateDb(res.data)
            }
        });

        function UpdateDb(filename) {
            axios.post("/profile_update_img", qs.stringify({
                Dp: filename,

            }), {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            });
        }
    } else {
        return;
    }
}

profile_edit_username.addEventListener("input", () => {
    if (profile_edit_username.value.length >= 3) {
        auth_errs.innerText = "";
        auth_errs.style.display = "none";
        axios.post("/check_user", qs.stringify({
            Username: profile_edit_username.value.trim(),
        }), {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }).then((res) => {
            if (res.data.Length > 0) {
                error("Username already taken");
            } else {
                axios.post("/profile_change_username", qs.stringify({
                    Username: profile_edit_username.value.trim(),

                }), {
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded'
                    }
                });
            }
        });
    } else {
        error("Username is too short");
    }
});

//on email input
profile_edit_email.addEventListener("input", () => {
    if (/^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(profile_edit_email.value)) {
        auth_errs.innerText = "";
        auth_errs.style.display = "none";
        axios.post("/check_user", qs.stringify({
            Email: profile_edit_email.value.trim(),
        }), {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }).then((res) => {
            if (res.data.Length > 0) {
                error("Email already in use");
            } else {
                axios.post("/profile_change_email", qs.stringify({
                    Email: profile_edit_email.value.trim(),
                }), {
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded'
                    }
                });
            }
        });
    } else {
        error("Email is invalid");
    }
});

//on bio input
profile_edit_bio.addEventListener("input", () => {
    auth_errs.innerText = "";
    auth_errs.style.display = "none";
    if (profile_edit_bio.value.trim().length > 0 && profile_edit_bio.value.trim().length < 210) {
        axios.post("/profile_change_bio", qs.stringify({
            Bio: profile_edit_bio.value.trim(),

        }), {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        });
    }
});