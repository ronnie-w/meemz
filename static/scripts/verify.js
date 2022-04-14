var qs = Qs;

const wrong_mail = document.getElementById("wrong_mail");
const verify_now = document.getElementById("verify_now");
const verify_code = document.getElementById("verify_code");
const err_msg = document.getElementById("err_msg");
const verify_msg = document.getElementById("verify_msg");

wrong_mail.addEventListener("click", () => {
    axios.post('/wrong_mail');
});

axios.post('/fetch_user').then(res => {
    verify_msg.innerText = `${res.data.Email}`;
});

verify_now.addEventListener("click", () => {
    axios.post('/verify_auth', qs.stringify({
        VCode: verify_code.value.trim()
    }), {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }).then((res) => {
        if (res.data.Err === undefined) {
            window.location.assign("/");
        }
        if (res.data.Err !== undefined) {
            err_msg.style.display = "block";
            err_msg.innerText = res.data.Err;
            err_msg.style.animation = "headShake";
            err_msg.style.animationDuration = "800ms";
        }
    }).catch((err) => {
        console.log(err);
    });
});