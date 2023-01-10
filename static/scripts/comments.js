import post_form from "./modules/dialer.js";
import { default as error_dialog, loading_dialog } from "./modules/dialog.js";
import comment_component from "./components/comment.js";

const meemz_post_comment = document.getElementById("meemz_post_comment"),
    meemz_post_btn = document.getElementById("meemz_post_btn"),
    url_parameters = window.location.pathname,
    file_id = url_parameters.split("/")[2];

let start_index = 0; 

(_ => {
    post_form("/fetch_comments", {
        FileId: file_id
    }).then(res => {
        res.data = res.data.reverse();
        function DisplayComments() {
            let comment = comment_component(res.data[start_index], "/delete_comment", "comment"),
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

            if (start_index < res.data.length-1) {
                start_index++;
                observer.observe(comment);
            }
        }
        
        DisplayComments();
    });
})()

meemz_post_btn.addEventListener("click", () => {
    if (meemz_post_comment.value.trim().length > 1) {
        let dialog = loading_dialog();
        post_form("/post_comment", {
            Comment: meemz_post_comment.value.trim(),
            FileId: file_id
        }).then(res => {
            if (res.data.Updated) {
                meemz_post_comment.value = "";
                dialog.remove();
                window.location.reload();
            } else if (res.data.Err != undefined) {
                error_dialog(res.data.Err, "/login");
            }
        });
    }
});