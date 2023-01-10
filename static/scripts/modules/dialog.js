import { default as post_form, upload, upload_config } from "./dialer.js";
import { default as button_loader, tags_processor } from "./loader.js";
import { default as element_constructor, _ } from "../components/element.js";

import comment_component from "../components/comment.js";
import multiple_slider from "./actions.js";
import theme_toggle from "../dependencies/theme.js";

const error_dialog = (error_message, url) => {
    if (url == undefined) {
        url = window.location.pathname;
    }

    document.querySelector("body").innerHTML +=
        `<dialog id="meemz_error_dialog" style="color: red;">
        <header>Error</header>
        <form method="dialog">
            <p>${error_message}</p>
            <center>
                <button onclick="javascript: window.location.assign('${url}')"><i class="fi fi-rr-cross-circle"></i></button>
            </center>
        </form>
    </dialog>`
        ;

    document.getElementById("meemz_error_dialog").showModal();
}

const media_dialog = (media_file, src, media_file_name, media_type, media_size, index, media_action) => {
    function dialog_popup(media, upload_url, config_url, file, filename) {
        let dialog = element_constructor("dialog", _, "meemz_media_dialog"),
            header = element_constructor("header", _, _, `${media_file_name}`),
            media_div = element_constructor("div", "meemz_media_div"),
            upload_div = element_constructor("div", _, "meemz_upload_div"),
            form = element_constructor("form"),
            media_upload_p = element_constructor("p"),
            media_size_p = element_constructor("p", _, _, `${round(media_size / 1000000)}mb`),
            media_actions = element_constructor("div", "meemz_media_actions"),
            caption_input = element_constructor("input"),
            tags_input = element_constructor("input"),
            credits_p = element_constructor("p"),
            credits_input = element_constructor("input"),
            cancel_btn = element_constructor("button"),
            upload_btn = element_constructor("button"),
            body = document.querySelector("body");

        media.style.maxHeight = `${screen.height - 200}px`;
        media.style.borderRadius = "6px";
        media.setAttribute("src", src);

        media_div.append(media);

        function round(float) {
            let m = Number((Math.abs(float) * 100).toPrecision(15));
            return Math.round(m) / 100 * Math.sign(float);
        }

        if (media_action !== undefined) {
            media_actions.innerHTML += media_action;
        } else {
            let media = document.getElementById(media_file_name);

            caption_input.setAttribute("placeholder", "Caption");
            tags_input.setAttribute("placeholder", "#tag1 #tag2 #tag3");
            credits_input.setAttribute("placeholder", "Credit a creator");

            credits_p.innerText = "Add a link or username to credit a creator ";
            credits_p.append(credits_input);

            cancel_btn.innerHTML = '<i class="fi fi-rr-cross-circle"></i>';
            upload_btn.innerHTML = '<i class="fi fi-rr-folder-upload"></i>';

            media_actions.append(cancel_btn, upload_btn);

            cancel_btn.addEventListener("click", () => {
                dialog.remove();
            });

            upload_btn.addEventListener("click", (e) => {
                e.preventDefault();

                cancel_btn.setAttribute("disabled", "disabled");

                media_upload_p.innerText = "Uploading...";

                let upload_res,
                    upload_config_res;

                button_loader(upload_btn);

                upload_res = upload(upload_url, file, filename);

                upload_res.then(res => {
                    if (res.data.Name !== '') {
                        upload_config_res = upload_config(res.data.Name, config_url, caption_input.value.trim(), tags_input.value.trim().toLocaleLowerCase(), credits_input.value.trim(), media_file_name, "single", index);

                        upload_config_res.then(res => {
                            if (res.data.Name === "Upload complete") {
                                media.remove();
                                dialog.remove();
                            }
                        });
                    }
                });

            });
        }

        form.setAttribute("method", "dialog");
        form.append(media_upload_p, media_size_p, caption_input, tags_input, credits_p, media_actions);

        dialog.append(header, media_div, upload_div, form);

        body.append(dialog);

        dialog.showModal();
    }

    if (media_type.includes("video")) {
        let video = document.createElement("video");
        video.setAttribute("controls", "controls");
        video.setAttribute("autoplay", "autoplay");
        video.setAttribute("loop", "loop");
        dialog_popup(video, "/upload_veemz", "/update_veemz_config", media_file, "veemz_upload");
    } else {
        let img = document.createElement("img");
        dialog_popup(img, "/upload_meemz", "/update_meemz_config", media_file, "meemz_upload");
    }
}

const loading_dialog = () => {
    let dialog = element_constructor("dialog"),
        center = element_constructor("center"),
        i = element_constructor("i", "fi fi-rr-spinner fa-spin meemz_loading_dialog"),
        body = document.querySelector("body");

    center.append(i);
    dialog.append(center);
    body.append(dialog);

    dialog.showModal();

    return dialog;
}

const plain_dialog = (message) => {
    let dialog = document.createElement("dialog"),
        message_p = element_constructor("p", _, _, message),
        center = element_constructor("center", _, "meemz_center_close"),
        button = document.createElement("button"),
        close = element_constructor("i", "fi fi-rr-cross-circle"),
        body = document.querySelector("body");

    button.append(close);
    center.append(button);
    dialog.append(message_p, center);

    button.addEventListener("click", () => {
        dialog.remove();
    });

    body.append(dialog);

    dialog.showModal();

    return dialog;
}

const profile_media_dialog = (media_files) => {
    let dialog = document.createElement("dialog"),
        media_div = document.createElement("div"),
        cancel_btn = document.createElement("button"),
        delete_btn = document.createElement("button"),
        confirm_delete_btn = document.createElement("button"),
        delete_actions = document.createElement("div"),
        share_btn = document.createElement("button"),
        media_img = document.createElement("img"),
        media_img_css = media_img.style,
        media_video = document.createElement("video"),
        media_caption = document.createElement("label"),
        media_credits_div = document.createElement("div"),
        media_credits = document.createElement("label"),
        media_tags = document.createElement("div"),
        media_actions = document.createElement("div"),
        body = document.querySelector("body");

    media_actions.setAttribute("class", "meemz_media_actions");
    media_actions.style.marginTop = "12px";

    delete_actions.setAttribute("class", "meemz_media_actions");
    delete_actions.style.marginTop = "12px";

    cancel_btn.innerHTML = '<i class="fi fi-rr-cross-circle"></i>';
    cancel_btn.addEventListener("click", () => {
        dialog.remove();
    });

    delete_btn.innerHTML = '<i class="fi fi-rr-trash"></i>';
    confirm_delete_btn.innerHTML = '<label>OK</label>';

    share_btn.innerHTML = '<i class="fi fi-rr-share"></i>';

    for (let i = 0; i < media_files.length; i++) {
        let file = media_files[i];
        var filename = file.FileName;

        media_credits_div.style.backgroundColor = "#0096bfab";
        media_credits_div.style.borderRadius = "6px";
        media_credits_div.style.border = "1px solid grey";
        media_credits_div.style.margin = "6px 0 6px 0";
        media_credits_div.style.width = "fit-content";

        media_credits.style.margin = "5px";

        file.Pcomment != "" ?
            media_caption.innerText = `${file.Pcomment}`
            : media_caption = "";

        file.Tags != "" ?
            media_tags = tags_processor(file.Tags)
            : media_tags = "";

        file.Credits != "" ?
            (_ => {
                media_credits.innerHTML = `${file.Credits}`;
                media_credits_div.append(media_credits);
            })()
            : media_credits_div = "";

        if (media_files.length > 1) {
            media_div = multiple_slider(media_files);
            break;
        } else {
            (_ => {
                filename.includes("veemz") ?
                    (_ => {
                        media_video.style.borderRadius = "6px";
                        media_video.style.maxHeight = `${screen.height - 300}px`;

                        media_video.setAttribute("src", `/static/veemz_uploads/${filename}`);
                        media_video.setAttribute("controls", "controls");
                        media_video.setAttribute("loop", "loop");
                        media_video.setAttribute("autoplay", "autoplay");
                        media_div.append(media_video);
                    })()
                    :
                    (_ => {
                        media_img_css.maxHeight = `${screen.height - 300}px`;
                        media_img_css.width = "100%";
                        media_img_css.borderRadius = "6px";
                        media_img_css.objectFit = "contain";
                        media_img_css.objectPosition = "center";

                        media_img.setAttribute("src", `/static/meemz_uploads/${filename}`);

                        media_div.append(media_img);
                    })()
            })()
        }
    }

    let share_content = {
        title: "meemz",
        text: media_caption.innerText,
        url: `/share/${media_files[0].FileId}`
    }

    share_btn.addEventListener("click", async () => {
        try {
            await navigator.share(share_content);
        } catch (err) {
            console.log(err);
        }
    });

    delete_btn.addEventListener("click", () => {
        let confirmation_dialog;
        if (media_files.length > 1) {
            confirmation_dialog = plain_dialog(`${media_files.length} files will be deleted. \nThis action is irreversible`);
        } else {
            confirmation_dialog = plain_dialog("Do you wish to delete this file? \nThis action is irreversible");
        }

        let center_btn_div = document.getElementById("meemz_center_close");
        center_btn_div.remove();

        cancel_btn.addEventListener("click", () => {
            confirmation_dialog.remove();
        });

        delete_actions.append(cancel_btn, confirm_delete_btn);
        confirmation_dialog.append(delete_actions);

        confirm_delete_btn.addEventListener("click", () => {
            cancel_btn.setAttribute("disabled", "disabled");
            button_loader(confirm_delete_btn);

            post_form("/delete_post", {
                FileId: media_files[0].FileId
            }).then(res => {
                res.data.Err != undefined ?
                    error_dialog(res.data.Err)
                    : (_ => {
                        let meemz_file_display = document.getElementsByClassName(`meemz_file_display ${filename}`)[0];
                        confirmation_dialog.remove();
                        dialog.remove();
                        meemz_file_display.remove();
                    })();
            });
        });
    });

    media_actions.append(cancel_btn, share_btn, delete_btn);

    dialog.append(media_div, media_caption, media_tags, media_credits_div, media_actions);

    body.append(dialog);

    dialog.showModal();
}

const reply_dialog = (comment_id) => {
    let dialog = document.createElement("dialog"),
        close_icon = document.createElement("i"),
        meemz_post_comment_div = document.createElement("div"),
        meemz_post_replies = document.createElement("div"),
        meemz_post_comment = document.createElement("textarea"),
        meemz_post_btn = document.createElement("button"),
        url_parameters = window.location.pathname,
        file_id = url_parameters.split("/")[2],
        start_index = 0,
        body = document.querySelector("body");

    (_ => {
        post_form("/fetch_replies", {
            FileId: file_id,
            CommentId: comment_id
        }).then(res => {
            if (res.data !== null) {
                res.data = res.data.reverse();
                function DisplayComments() {
                    let comment = comment_component(res.data[start_index], "/delete_reply", "reply"),
                        observer = new IntersectionObserver(
                            entries => {
                                entries.map(entry => {
                                    if (entry.isIntersecting) {
                                        DisplayComments();
                                        observer.unobserve(comment);
                                    }
                                });
                            }, {
                            threshold: 1.0
                        }
                        );

                    if (start_index < res.data.length - 1) {
                        start_index++;
                        observer.observe(comment);
                    }

                    meemz_post_replies.append(comment);
                }

                DisplayComments(res.data);
            }
        });
    })()

    meemz_post_btn.innerHTML = `<i class="fi fi-rr-interface"></i>`;

    meemz_post_comment.setAttribute("rows", "2");
    meemz_post_comment.setAttribute("placeholder", "Reply to this comment");

    meemz_post_comment_div.append(meemz_post_comment, meemz_post_btn);

    close_icon.addEventListener("click", () => {
        dialog.remove();
    });

    close_icon.setAttribute("class", "fi fi-rr-cross-circle");
    close_icon.style.marginBottom = "6px";
    close_icon.style.float = "right";

    meemz_post_btn.addEventListener("click", () => {
        if (meemz_post_comment.value.length > 0) {
            post_form("/post_reply", {
                Comment: meemz_post_comment.value.trim(),
                FileId: file_id,
                CommentId: comment_id
            }).then(res => {
                if (res.data.Updated) {
                    dialog.remove();
                } else {
                    error_dialog(res.data.Err, "/login");
                }
            });
        }
    });

    dialog.append(close_icon, meemz_post_comment_div, meemz_post_replies);

    body.append(dialog);

    dialog.showModal();
}

const theme_dialog = () => {
    let dialog = document.createElement("dialog"),
        close_btn = element_constructor("button", _, _, _, _, '<i class="fi fi-rr-cross-circle"></i>'),
        media_actions = document.createElement("center"),
        header = element_constructor("header", _, _, "Choose a theme"),
        light_div = element_constructor("div", "meemz_change_theme"),
        dark_div = element_constructor("div", "meemz_change_theme"),
        auto_div = element_constructor("div", "meemz_change_theme"),
        light_p = element_constructor("p", _, _, "Light"),
        dark_p = element_constructor("p", _, _, "Dark"),
        auto_p = element_constructor("p", _, _, "Automatic"),
        body = document.querySelector("body");

    media_actions.style.width = "100%";
    media_actions.style.marginTop = "12px";

    close_btn.style.marginTop = "12px";
    close_btn.addEventListener("click", () => {
        dialog.remove();
    });

    light_div.addEventListener("click", () => {
        theme_toggle("light");
        window.location.reload();
    });

    dark_div.addEventListener("click", () => {
        theme_toggle("dark");
        window.location.reload();
    });

    auto_div.addEventListener("click", () => {
        theme_toggle("auto");
        window.location.reload();
    });

    light_div.append(light_p);
    dark_div.append(dark_p);
    auto_div.append(auto_p);

    media_actions.append(close_btn);
    dialog.append(header, auto_div, light_div, dark_div, media_actions);

    body.append(dialog);

    dialog.showModal();
}

const profile_edit_dialog = (url, placeholder) => {
    let dialog = element_constructor("dialog"),
        input = element_constructor("input"),
        button = element_constructor("button", _, _, "Change"),
        close_btn = element_constructor("button", _, _, _, _, '<i class="fi fi-rr-cross-circle"></i>'),
        actions_div = element_constructor("div", "meemz_media_actions"),
        body = document.querySelector("body");

    input.setAttribute("placeholder", placeholder);

    close_btn.addEventListener("click", () => {
        dialog.remove();
    });

    button.addEventListener("click", () => {
        button_loader(button);
        post_form(url, { Data: input.value.trim() }).then(res => {
            if (res.data.Err !== undefined) {
                error_dialog(res.data.Err);
            } else {
                window.location.reload();
            }
        })
    });

    actions_div.append(button, close_btn);
    dialog.append(input, actions_div);
    body.append(dialog);

    dialog.showModal();
}

export { error_dialog as default, media_dialog, loading_dialog, plain_dialog, profile_media_dialog, reply_dialog, theme_dialog, profile_edit_dialog };
