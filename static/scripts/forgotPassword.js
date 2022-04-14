var qs = Qs;

var current_state;

const username_div = document.getElementById("username_div");
const verify_div = document.getElementById("verify_div");
const newpass_div = document.getElementById("newpass_div");
const verify_msg = document.getElementById("verify_msg");
const username = document.getElementById("username");
const v_code = document.getElementById("v_code");
const password = document.getElementById("password");
const next_btn = document.getElementById("next_btn");
const user_err = document.getElementById("user_err");
const verify_err = document.getElementById("verify_err");
const pass_err = document.getElementById("pass_err");

function CurrentState() {
    console.log(current_state);
    if (current_state === undefined) {
        next_btn.innerHTML = `<i style="color : black;" class="fas fa-spinner-third fa-spin"></i>`;
        axios.post("/pass_send_code", qs.stringify({
            Username: username.value.trim()
        }), {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }).then(res => {
            next_btn.innerHTML = "Next";
            if (res.data.Err === undefined) {
                verify_msg.innerText = `A verification code has been sent to ${res.data.Email}`;
                current_state = "verify";
                username_div.style.display = "none";
                verify_div.style.display = "grid";
                newpass_div.style.display = "none";
            } else {
                user_err.innerText = `${res.data.Err}`;
            }
        });
    } else if (current_state === "verify") {
        next_btn.innerHTML = `<i style="color : black;" class="fas fa-spinner-third fa-spin"></i>`;
        user_err.innerText = "";
        //let auth_key = document.cookie.match(new RegExp('(^| )uid=([^;]+)'))[2];
        axios.post("/pass_verify_code", qs.stringify({
            PassCode: v_code.value.trim()
        }), {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }).then(res => {
            next_btn.innerHTML = "Next";
            if (res.data.Err !== undefined) {
                verify_err.innerText = `${res.data.Err}`;
            } else {
                current_state = "new_password";
                username_div.style.display = "none";
                verify_div.style.display = "none";
                newpass_div.style.display = "grid";
            }
        });
    } else if (current_state === "new_password") {
        verify_err.innerText = "";
        if (password.value.trim().length < 6) {
            pass_err.innerText = `Password length should be 6 or more`
        } else {
           // let auth_key = document.cookie.match(new RegExp('(^| )uid=([^;]+)'))[2];
            axios.post("/pass_new_password", qs.stringify({
                NewPassword: password.value.trim()
            }), {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            }).then(_res => {
                window.location.replace("/login");
            });
        }
    }
}