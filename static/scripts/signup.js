var qs = Qs;

const username_input = document.getElementById("username_input");
const email_input = document.getElementById("email_input");
const password_input = document.getElementById("password_input");
const show_hide = document.getElementById("show_hide");
const auth_errs = document.getElementById("auth_errs");
const signup_btn = document.getElementById("signup_btn");

function error(msg) {
    auth_errs.style.display = "block";
    auth_errs.innerText = msg;
    auth_errs.style.animation = "headShake";
    auth_errs.style.animationDuration = "800ms";
}

username_input.addEventListener("input", () => {
    auth_errs.innerText = "";
    auth_errs.style.display = "none";
    axios.post("/check_user", qs.stringify({
        Username: username_input.value.trim()
    }), {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }).then((res) => {
        if (res.data.Length > 0) {
            error("Username already taken");
        }
    });
});

email_input.addEventListener("input", () => {
    auth_errs.innerText = "";
    auth_errs.style.display = "none";
    axios.post("/check_user", qs.stringify({
        Email: email_input.value.trim()
    }), {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }).then((res) => {
        if (res.data.Length > 0) {
            error("Email already in use");
        }
    });
});

password_input.addEventListener("input", () => {
    auth_errs.innerText = "";
    auth_errs.style.display = "none";
    if (password_input.value.length < 6) {
        error("Password length must be six or more");
    }
});

show_hide.addEventListener("click", () => {
    switch (password_input.getAttribute('type')) {
        case 'password':
            show_hide.setAttribute('class', 'password_eye fal fa-eye-slash')
            password_input.setAttribute('type', 'text');
            break;

        case 'text':
            show_hide.setAttribute('class', 'password_eye fal fa-eye')
            password_input.setAttribute('type', 'password');
            break;

        default:
            break;
    }
});

function submitSignup() {
    var non_regex = /(!|@|#|%|&|-|=|:|"|}|{|]|'|;|<|>|,|~|`)/gi;
    var matcher = /[a-z]/gi;

    if (username_input.value.length < 3) {
        error("Username is too short");
    } else if (username_input.value.match(non_regex) || !username_input.value.match(matcher)) {
        error("Invalid username");
    } else if (password_input.value.length < 6) {
        error("Password is too short");
    } else if (!/^\w+([.-]?\w+)*@\w+([.-]?\w+)*(\.\w{2,3})+$/.test(email_input.value)) {
        error("Invalid email address");
    } else {
        signup_btn.setAttribute("disabled", "disabled");
        signup_btn.innerHTML = `<i style="color : #121212;" class="fas fa-spinner-third fa-spin"></i>`;
        axios.post('/signup_auth', qs.stringify({
            Username: username_input.value.trim(),
            Email: email_input.value.trim(),
            Password: password_input.value.trim()
        }), {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }).then((res) => {
            if (res.data.Err === undefined) {
                window.location.assign("/verify");
            }
            if (res.data.Err !== undefined) {
                error(res.data.Err);
            }
        }).catch((err) => {
            console.log(err);
        });
    }
}