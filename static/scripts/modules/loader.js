'use strict'

const button_loader = (btn) => {
    btn.innerHTML = `<i class="fi fi-rr-spinner fa-spin"></i>`;
    btn.setAttribute("disabled", "disabled");
}

const tags_processor = (tags) => {
    let tags_arr = tags.split(" "),
        tags_div = document.createElement("div");

    for (let i = 0; i < tags_arr.length; i++) {
        let tag,
            tag_div = document.createElement("div"),
            tag_label = document.createElement("label"),
            tag_div_css = tag_div.style,
            tag_label_css = tag_label.style;

        tags_arr[i].includes("#") ?
            tag = tags_arr[i].replace("#", "")
            : tag = tags_arr[i];

        tag_label_css.margin = "5px";

        tag_div_css.border = "1px solid grey";
        tag_div_css.borderRadius = "6px";
        tag_div_css.width = "fit-content";
        tag_div_css.display = "inline-block";
        tag_div_css.margin = "3px 6px 3px 0";

        tag_label.innerText = `${tag}`;

        tag_div.append(tag_label);

        tags_div.append(tag_div);
    }

    tags_div.style.margin = "12px 0 6px 0";

    return tags_div;
}

export { button_loader as default, tags_processor };