import post_form from "../modules/dialer.js";
import button_loader from "../modules/loader.js";
import { plain_dialog, reply_dialog } from "../modules/dialog.js";
import { default as element_constructor, _ } from "./element.js";

const comment_component = (comment, delete_url, comment_type) => {
    let meemz_comment = document.getElementById("meemz_comment"),
        meemz_comment_div = element_constructor("div", "meemz_comment_div"),
        meemz_comment_details = element_constructor("div", "meemz_comment_details"),
        meemz_comment_profile_img_div = element_constructor("div", "meemz_comment_profile_img_div"),
        meemz_comment_profile_img = element_constructor("img", "meemz_comment_profile_img", _, _, `/static/profile-pictures/${comment.ProfileImg}`),
        meemz_comment_username_div = element_constructor("div", "meemz_comment_username_div"),
        meemz_comment_username_label = element_constructor("label", _, _, `${comment.Username}`),
        meemz_comment_timestamp_div = element_constructor("div", "meemz_comment_timestamp_div", _, `${comment.CommentTime}`),
        meemz_comment_timestamp_label = element_constructor("label"),
        meemz_comment_center = element_constructor("center"),
        meemz_main_comment_div = element_constructor("div", "meemz_main_comment_div"),
        meemz_main_comment = element_constructor("label", _, _, _, _, `${comment.Comment}`),
        meemz_comment_actions = element_constructor("div", "meemz_comment_actions"),
        meemz_comment_like = element_constructor("div", "meemz_comment_like"),
        meemz_comment_like_action = element_constructor("i"),
        meemz_comment_like_count = element_constructor("div", "meemz_comment_like_count"),
        meemz_comment_count_num = element_constructor("label", "meemz_comment_count_num", _, `${comment.CommentLikes}`),
        meemz_comment_reply = element_constructor("div", "meemz_comment_reply"),
        meemz_comment_reply_action = element_constructor("i", "fi fi-rr-comments"),
        meemz_comment_reply_count = element_constructor("div", "meemz_comment_reply_count"),
        meemz_comment_count_reply = element_constructor("label", "meemz_comment_count_num", _, `${comment.Replies}`),
        meemz_comment_hide = element_constructor("div", "meemz_comment_hide"),
        meemz_comment_ban = element_constructor("i", "fi fi-rr-ban"),
        meemz_comment_delete = element_constructor("div", "meemz_comment_delete"),
        meemz_comment_delete_ico = element_constructor("i", "fi fi-rr-trash"),
        cancel_confirm = element_constructor("button"),
        confirmation_btn = element_constructor("button"),
        url_parameters = window.location.pathname,
        file_id = url_parameters.split("/")[2];
        
    if (comment.CommentLikedByUser) {
        meemz_comment_like_action.setAttribute("class", "fi fi-br-comment-heart");
    } else {
        meemz_comment_like_action.setAttribute("class", "fi fi-rr-comment-heart");
    }

    meemz_comment_details.append(meemz_comment_profile_img_div, meemz_comment_username_div, meemz_comment_timestamp_div);
    meemz_comment_profile_img_div.append(meemz_comment_profile_img);
    meemz_comment_username_div.append(meemz_comment_username_label);
    meemz_comment_timestamp_div.append(meemz_comment_timestamp_label);

    meemz_comment_center.append(meemz_main_comment_div);
    meemz_main_comment_div.append(meemz_main_comment);

    meemz_comment_like.append(meemz_comment_like_action, meemz_comment_like_count);
    meemz_comment_like_count.append(meemz_comment_count_num);
    meemz_comment_reply.append(meemz_comment_reply_action, meemz_comment_reply_count);
    meemz_comment_reply_count.append(meemz_comment_count_reply);
    meemz_comment_hide.append(meemz_comment_ban);
    meemz_comment_delete.append(meemz_comment_delete_ico);

    meemz_comment_reply.addEventListener("click", () => {
        reply_dialog(comment.CommentId);
    });

    function DeleteDialog(message) {
        let confirmation_dialog = plain_dialog(message),
            cancel_confirmation = document.getElementById("meemz_center_close"),
            media_actions = element_constructor("div");

        confirmation_btn.innerText = "OK";

        cancel_confirmation.remove();

        cancel_confirm.innerHTML = `<i class="fi fi-rr-cross-circle"></i>`;

        media_actions.append(cancel_confirm, confirmation_btn);
        confirmation_dialog.append(media_actions);

        media_actions.setAttribute("class", "meemz_media_actions");
        media_actions.style.marginTop = "12px";

        cancel_confirm.addEventListener("click", () => {
            confirmation_dialog.remove();
        });

        confirmation_btn.addEventListener("click", () => {
            cancel_confirm.setAttribute("disabled", "disabled");
            button_loader(confirmation_btn);

            post_form(delete_url, {
                CommentId: comment.CommentId
            }).then(res => {
                res.data.Updated ?
                    (_ => {
                        confirmation_dialog.remove();
                        meemz_comment_div.remove();
                        console.log(comment);
                    })() : null;
            });
        });
    }

    function HideDialog(message) {
        let confirmation_dialog = plain_dialog(message),
            cancel_confirmation = document.getElementById("meemz_center_close"),
            media_actions = element_constructor("div");

        confirmation_btn.innerText = "OK";

        cancel_confirmation.remove();

        cancel_confirm.innerHTML = `<i class="fi fi-rr-cross-circle"></i>`;

        media_actions.append(cancel_confirm, confirmation_btn);
        confirmation_dialog.append(media_actions);

        media_actions.setAttribute("class", "meemz_media_actions");
        media_actions.style.marginTop = "12px";

        cancel_confirm.addEventListener("click", () => {
            confirmation_dialog.remove();
        });

        confirmation_btn.addEventListener("click", () => {
            cancel_confirm.setAttribute("disabled", "disabled");
            button_loader(confirmation_btn);

            confirmation_dialog.remove();
            meemz_comment_div.remove();
        });
    }

    meemz_comment_delete_ico.addEventListener("click", () => {
        DeleteDialog(`This ${comment_type} will be deleted`)
    });

    meemz_comment_ban.addEventListener("click", () => {
        HideDialog(`This ${comment_type} will be hidden`);
    });

    meemz_comment_like_action.addEventListener("click", () => {
        let like_count = meemz_comment_count_num.innerText, current_count, class_name = meemz_comment_like_action.getAttribute("class"),
            update_count = (parse_class, count, color) => {
                meemz_comment_like_action.setAttribute("class", parse_class);
                meemz_comment_like_action.style.color = color;
                current_count = parseInt(like_count);
                meemz_comment_count_num.innerText = `${current_count + count}`;
            }
        class_name === "fi fi-rr-comment-heart" ?
            (_ => {
                update_count("fi fi-br-comment-heart animate__animated animate__rubberBand", 1, "red");
                post_form("/post_comment_reaction", {
                    FileId: file_id,
                    CommentId: comment.CommentId
                });
            })() : (_ => {
                update_count("fi fi-rr-comment-heart", -1, "");
                post_form("/delete_comment_reaction", {
                    CommentId: comment.CommentId
                });
            })()
    });


    delete_url === "/delete_comment" ?
        (_ => {
            meemz_comment.append(meemz_comment_div)
            if (comment.CommentAction) {
                meemz_comment_actions.append(meemz_comment_like, meemz_comment_reply, meemz_comment_hide, meemz_comment_delete);
            } else {
                meemz_comment_actions.append(meemz_comment_like, meemz_comment_reply, meemz_comment_hide);
            }
        })()
        : (_ => {
            if (comment.CommentAction) {
                meemz_comment_actions.append(meemz_comment_like, meemz_comment_hide, meemz_comment_delete);
            } else {
                meemz_comment_actions.append(meemz_comment_like, meemz_comment_hide);
            }
        })()

    meemz_comment_div.append(meemz_comment_details, meemz_comment_center, meemz_comment_actions);
    return meemz_comment_div;
}

export default comment_component;