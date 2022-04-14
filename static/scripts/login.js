var qs = Qs;

const login_username = document.getElementById("login_username");
const login_password = document.getElementById("login_password");
const login_btn = document.getElementById("login_btn");
const auth_errs = document.getElementById("auth_errs");
const show_hide = document.getElementById("show_hide");

function error(msg) {
    auth_errs.style.display = "block";
    auth_errs.innerText = msg;
    auth_errs.style.animation = "headShake";
    auth_errs.style.animationDuration = "800ms";
}

login_btn.addEventListener("click", () => {
    axios.post('/login_auth', qs.stringify({
        Username: login_username.value.trim(),
        Password: login_password.value.trim()
    }), {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }).then((res) => {
        if (res.data.Err === undefined) {
            window.location.assign("/");
        }
        if (res.data.Err !== undefined) {
            error(res.data.Err);
        }
    }).catch((err) => {
        console.log(err);
    });
});

show_hide.addEventListener("click", () => {
    switch (login_password.getAttribute('type')) {
        case 'password':
            show_hide.setAttribute('class', 'password_eye fal fa-eye-slash')
            login_password.setAttribute('type', 'text');
            break;

        case 'text':
            show_hide.setAttribute('class', 'password_eye fal fa-eye')
            login_password.setAttribute('type', 'password');
            break;

        default:
            break;
    }
});