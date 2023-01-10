import post_form from "./modules/dialer.js";
import button_loader from "./modules/loader.js";
import error_dialog from "./modules/dialog.js";

const login = document.querySelector("body > center:nth-child(2) > div > div:nth-child(1)"),
    signup = document.querySelector("body > center:nth-child(2) > div > div:nth-child(2)"),
    auth_btn = document.querySelector("body > center:nth-child(4) > button"),
    show_password = document.querySelector("#show_password"),
    path = window.location.pathname;

let username = document.getElementById("meemz_username"),
    password = document.getElementById("meemz_password"),
    email, verification_code;

function show_hide() {
    show_password.addEventListener("click", () => {
        if (show_password.checked == true) {
            password.removeAttribute("type");
            password.setAttribute("type", "text");
        } else {
            password.removeAttribute("type");
            password.setAttribute("type", "password");
        }
    });
}

path == "/login" || path == "/signup"
    ? login.addEventListener("click", () => {
        window.location.assign("/login");
    }) || signup.addEventListener("click", () => {
        window.location.assign("/signup");
    }) : null;

path == "/login"
    ? (_ => {
        show_hide();
        login.style.backgroundColor = "#0096bfab";
        auth_btn.addEventListener("click", () => {
            button_loader(auth_btn);
            post_form("/login_auth", {
                Username: username.value.trim(),
                Password: password.value.trim()
            }).then(res => {
                if (res.data.Err != undefined) {
                    error_dialog(res.data.Err);
                } else {
                    window.location.assign("/");
                }
            });
        });
    })()
    : path == "/signup"
        ? (_ => {
            show_hide();
            signup.style.backgroundColor = "#0096bfab";
            auth_btn.addEventListener("click", () => {
                button_loader(auth_btn);
                email = document.getElementById("meemz_email");
                post_form("/signup_auth", {
                    Username: username.value.trim(),
                    Email: email.value.trim(),
                    Password: password.value.trim()
                }).then(res => {
                    if (res.data.Err != undefined) {
                        error_dialog(res.data.Err);
                    } else {
                        window.location.assign("/verify");
                    }
                });
            });
        })()
        : path == "/verify"
            ? (_ => {
                auth_btn.addEventListener("click", () => {
                    button_loader(auth_btn);
                    verification_code = document.getElementById("meemz_verification_code");
                    post_form("/verify_auth", {
                        VCode: verification_code.value.trim()
                    }).then(res => {
                        if (res.data.Err != undefined) {
                            error_dialog(res.data.Err);
                        } else {
                            window.location.replace("/");
                        }
                    });
                });
            })()
            : path == "/forgot_password"
                ? (_ => {
                    auth_btn.addEventListener("click", () => {
                        button_loader(auth_btn);
                        post_form("/pass_send_code", {
                            Username: username.value.trim()
                        }).then(res => {
                            if (res.data.Err != undefined) {
                                error_dialog(res.data.Err);
                            } else {
                                window.location.replace("/verify_passcode");
                            }
                        });
                    });
                })()
                : path == "/verify_passcode"
                    ? (_ => {
                        auth_btn.addEventListener("click", () => {
                            button_loader(auth_btn);
                            verification_code = document.getElementById("meemz_verification_code");
                            post_form("/pass_verify_code", {
                                PassCode: verification_code.value.trim()
                            }).then(res => {
                                if (res.data.Err != undefined) {
                                    error_dialog(res.data.Err);
                                } else {
                                    window.location.replace("/password_reset");
                                }
                            });
                        });
                    })()
                    : path == "/password_reset"
                        ? (_ => {
                            show_hide();
                            auth_btn.addEventListener("click", () => {
                                button_loader(auth_btn);
                                post_form("/pass_new_password", {
                                    NewPassword: password.value.trim()
                                }).then(res => {
                                    if (res.data.Err != undefined) {
                                        error_dialog(res.data.Err);
                                    } else {
                                        window.location.replace("/login");
                                    }
                                });
                            });
                        })() : null;



