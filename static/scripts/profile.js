import { default as post_form, get_req, upload } from "./modules/dialer.js";
import { default as error_dialog, profile_media_dialog, plain_dialog, theme_dialog, profile_edit_dialog } from "./modules/dialog.js";
import button_loader from "./modules/loader.js";

const meemz_profile_pic = document.getElementById("meemz_profile_pic"),
    meemz_profile_username = document.getElementById("meemz_profile_username"),
    meemz_profile_bio = document.getElementById("meemz_profile_bio"),
    meemz_profile_uploads = document.getElementById("meemz_profile_uploads"),
    meemz_folder = document.getElementById("meemz_folder"),
    meemz_settings = document.getElementById("meemz_settings"),
    meemz_profile_uploads_div = document.getElementById("meemz_profile_uploads_div"),
    meemz_profile_settings_div = document.getElementById("meemz_profile_settings_div"),
    meemz_edit_dp = document.getElementById("meemz_edit_dp"),
    meemz_dp_select = document.getElementById("meemz_dp_select"),
    meemz_profile_details = document.getElementById("meemz_profile_details"),
    meemz_display_username = document.getElementById("meemz_display_username"),
    meemz_display_email = document.getElementById("meemz_display_email"),
    meemz_display_bio = document.getElementById("meemz_display_bio"),
    meemz_display_theme = document.getElementById("meemz_display_theme"),
    meemz_edit_username_div = document.getElementById("meemz_edit_username_div"),
    meemz_edit_email_div = document.getElementById("meemz_edit_email_div"),
    meemz_edit_bio_div = document.getElementById("meemz_edit_bio_div"),
    meemz_theme = document.getElementById("meemz_theme"),
    meemz_advertising = document.getElementById("meemz_advertising"),
    meemz_logout = document.getElementById("meemz_logout");

let uploads_arr, start_index = 0, end_index;

function sort(mediafiles) {
    let sorted = [];
    for (let i = 0; i < mediafiles.length; i++) {
        sorted[mediafiles[i].FileIndex] = mediafiles[i];
    }

    return sorted;
}

(_ => {
    get_req("/fetch_user").then(res => {
        if (res.data.Err !== undefined) {
            error_dialog(res.data.Err, "/login");
        } else {
            meemz_profile_pic.style.backgroundImage = `url('/static/profile-pictures/${res.data.ProfileImg}')`;
            meemz_profile_username.innerHTML = `<label style="margin: 6px;">${res.data.Username}</label>`;

            res.data.Bio !== "No bio found" ?
                (_ => {
                    meemz_profile_bio.innerHTML = `<label style="margin: 6px;">${res.data.Bio}</label>`;
                    meemz_profile_details.style.position = "absolute";
                })() : null;

            meemz_display_username.innerText = `${res.data.Username}`;
            meemz_display_email.innerText = `${res.data.Email}`;
            meemz_display_bio.innerText = `${res.data.Bio}`;
            meemz_display_theme.innerText = `${document.cookie.match(new RegExp('(^| )color-scheme=([^;]+)'))[2]}`;
        }
    });

    function MediaDisplay(start_index, end_index, media) {
        for (let i = start_index; i < end_index; i++) {
            if (i == media.length) {
                break;
            }

            if (media[i].length > 1) {
                media[i] = sort(media[i]);
            }

            let file = media[i][0],
                filename = file.FileName,
                meemz_file_display = document.createElement("div"),
                observer = new IntersectionObserver(
                    entries => {
                        entries.map(entry => {
                            entry.isIntersecting && filename.includes("veemz") ?
                                meemz_file_display.innerHTML = "<i class='fi fi-br-play' style='position: relative; float: left; margin: 6px;'></i>"
                                :
                                entry.isIntersecting && filename.includes("meemz") ?
                                    meemz_file_display.style.backgroundImage = `url('/static/meemz_uploads/${filename}')`
                                    : null;

                            i == end_index - 1 ?
                                entry.isIntersecting ?
                                    (_ => {
                                        start_index = end_index;
                                        observer.unobserve(meemz_file_display);
                                        MediaDisplay(start_index, end_index + 3, media);
                                    })()
                                    :
                                    null : null;
                        });
                    },
                    {
                        threshold: 1.0
                    });

            meemz_file_display.setAttribute("class", `meemz_file_display ${filename}`);
            meemz_file_display.style.cursor = "pointer";

            meemz_profile_uploads.append(meemz_file_display);

            observer.observe(meemz_file_display);

            meemz_file_display.addEventListener("click", () => {
                profile_media_dialog(media[i]);
            });
        }
    }

    get_req("/my_uploads").then(res => {
        res.data != null ?
            (_ => {
                uploads_arr = res.data.reverse();
                MediaDisplay(start_index, start_index + 3, uploads_arr);
            })()
            :
            meemz_profile_uploads.innerText = "No posts available";
    });
})()

const display_change = (icon, display) => {
    let match = document.cookie.match(new RegExp('(^| )color-scheme=([^;]+)'))[2],
        dark = () => {
            meemz_folder.style.color = "#dbdbdb";
            meemz_settings.style.color = "#dbdbdb";
        },
        light = () => {
            meemz_folder.style.color = "#363636";
            meemz_settings.style.color = "#363636";
        };

    if (match == "dark") {
        dark();
    } else if (match == "light") {
        light();
    } else if (match == "auto" && window.matchMedia("(prefers-color-scheme: dark)").matches) {
        dark();
    } else if (match == "auto" && window.matchMedia("(prefers-color-scheme: light)").matches) {
        light();
    }
    
    icon.style.color = "#0096bfab";

    meemz_profile_uploads_div.style.display = "none";
    meemz_profile_settings_div.style.display = "none";

    display.style.display = "block";
}

// meemz_dp_edit.addEventListener("change", () => {
//     if (meemz_dp_edit.files.length > 0) {
//         let reader = new FileReader();
//         reader.onload = () => {
//             meemz_profile_pic.style.backgroundImage = `url('${reader.result}')`;
//         }
//         reader.readAsDataURL(meemz_dp_edit.files[0]);

//         meemz_dp_edit_div.innerHTML = `<label for="meemz_dp_edit" class="meemz_dp_edit">Image selected</label>`;
//     }
// });

meemz_folder.addEventListener("click", () => {
    display_change(meemz_folder, meemz_profile_uploads_div);
});

meemz_settings.addEventListener("click", () => {
    display_change(meemz_settings, meemz_profile_settings_div);
});

meemz_theme.addEventListener("click", () => {
    theme_dialog();
});

meemz_edit_username_div.addEventListener("click", () => {
    profile_edit_dialog("/profile_change_username", "Username");
});

meemz_edit_email_div.addEventListener("click", () => {
    profile_edit_dialog("/profile_change_email", "Email");
});

meemz_edit_bio_div.addEventListener("click", () => {
    profile_edit_dialog("/profile_change_bio", "Bio");
});

meemz_edit_dp.addEventListener("click", () => {
    meemz_dp_select.addEventListener("change", () => {
        if (meemz_dp_select.files.length > 0) {
            let reader = new FileReader();
            reader.onload = () => {
                meemz_profile_pic.style.backgroundImage = `url('${reader.result}')`;
            }

            reader.readAsDataURL(meemz_dp_select.files[0]);

            upload("/profile_img_upload", meemz_dp_select.files[0], "profile_img").then(res => {
                post_form("/profile_update_img", { Dp: res.data });
            });
        }
    });
});

meemz_advertising.addEventListener("click", () => {
    plain_dialog("Share this app to friends and family and make this possible :)");
});

meemz_logout.addEventListener("click", () => {
    document.cookie = "uid=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    window.location.replace("/login");
});