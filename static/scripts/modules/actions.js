const multiple_slider = (media_files) => {
    let media_div = document.createElement("div"),
        media = document.createElement("div"),
        label_div = document.createElement("div"),
        label = document.createElement("label"),
        img = document.createElement("img"),
        vid = document.createElement("video"),
        spinner = document.createElement("i"),
        index = 0;

    label_div.style.position = "absolute";
    label_div.style.margin = "6px";
    label_div.style.float = "left";
    label_div.style.color = "white";
    label_div.style.background = "grey";
    label_div.style.borderRadius = "6px";
    label_div.style.opacity = "0.5";

    label.style.margin = "6px";
    label_div.append(label);

    function sort(mediafiles) {
        let sorted = [];
        for (let i = 0; i < mediafiles.length; i++) {
            sorted[mediafiles[i].FileIndex] = mediafiles[i];
        }

        return sorted;
    }

    media_files = sort(media_files);

    media_files[0].FileName.includes("veemz") ?
        (_ => {
            vid.setAttribute("src", `/static/veemz_uploads/${media_files[0].FileName}`);
            vid.setAttribute("controls", "controls");
            vid.setAttribute("loop", "loop");
            vid.setAttribute("autoplay", "autoplay");

            label.innerText = `${index + 1}/${media_files.length}`;

            media.append(label_div, vid);
        })()
        :
        (_ => {
            img.setAttribute("src", `/static/meemz_uploads/${media_files[0].FileName}`);
            label.innerText = `${index + 1}/${media_files.length}`;

            img.addEventListener("load", () => {
                media.append(label_div, img);
            });
        })()

    media_div.append(media);

    img.style.width = "100%";
    img.style.maxHeight = `${screen.height - 300}px`;
    img.style.objectFit = "contain";
    img.style.objectPosition = "center";

    vid.style.maxHeight = `${screen.height - 300}px`;
    vid.style.objectFit = "contain";
    vid.style.objectPosition = "center";

    media_div.setAttribute("class", "meemz_media_multiple_div");

    let hammer = new Hammer(media);

    function video(src) {
        vid.setAttribute("src", `/static/veemz_uploads/${src}`);
        vid.setAttribute("controls", "controls");
        vid.setAttribute("loop", "loop");
        vid.setAttribute("autoplay", "autoplay");
    }

    hammer.on("panend", ev => {
        img.remove();
        vid.remove();
        label_div.remove();

        if (ev.additionalEvent == "panleft") {
            index < media_files.length - 1 ?
                index++ : null;
            console.log(ev.additionalEvent, media_files[index]);

            media_files[index].FileName.includes("veemz") ?
                (_ => {
                    video(media_files[index].FileName);
                    label.innerText = `${index + 1}/${media_files.length}`;
                    media.append(label_div, vid);
                })()
                :
                (_ => {
                    spinner.remove();
                    spinner.setAttribute("class", "fi fi-rr-spinner fa-spin meemz_loading_multiple");
                    media.style.height = "300px";
                    media.append(spinner);

                    img.setAttribute("src", `/static/meemz_uploads/${media_files[index].FileName}`);
                    img.addEventListener("load", _ => {
                        spinner.remove();
                        label.innerText = `${index + 1}/${media_files.length}`;
                        media.style.height = "fit-content";
                        media.append(label_div, img);
                    });
                })()
        } else {
            index > 0 ?
                index-- : null;
            console.log(ev.additionalEvent, media_files[index]);
            media_files[index].FileName.includes("veemz") ?
                (_ => {
                    video(media_files[index].FileName);
                    label.innerText = `${index + 1}/${media_files.length}`;
                    media.append(label_div, vid);
                })()
                :
                (_ => {
                    spinner.remove();
                    spinner.setAttribute("class", "fi fi-rr-spinner fa-spin meemz_loading_multiple");
                    media.style.height = "300px";
                    media.append(spinner);

                    img.setAttribute("src", `/static/meemz_uploads/${media_files[index].FileName}`);
                    img.addEventListener("load", _ => {
                        spinner.remove();
                        label.innerText = `${index + 1}/${media_files.length}`;
                        media.style.height = "fit-content";
                        media.append(label_div, img);
                    });
                })()
        }
    });

    return media_div;
}

export default multiple_slider;