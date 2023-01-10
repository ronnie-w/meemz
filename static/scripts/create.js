import { upload, upload_config, get_req } from "./modules/dialer.js";
import { default as error_dialog, media_dialog, plain_dialog, loading_dialog } from "./modules/dialog.js";
import button_loader from "./modules/loader.js";

const display_div = document.getElementById("meemz_files_selected"),
   upload_btn = document.getElementById("meemz_upload_btn"),
   upload_single = document.getElementById("meemz_upload_single_div"),
   upload_multiple = document.getElementById("meemz_upload_multiple_div"),
   upload_caption = document.getElementById("meemz_upload_caption"),
   upload_tags = document.getElementById("meemz_upload_tags"),
   upload_credits = document.getElementById("meemz_upload_credits"),
   upload_num = document.getElementById("meemz_selected_num");

var file_input = document.getElementById("meemz_select_file"),
   upload_type = "single";

const file_format = (file_element, file, files_length, i) => {
   let reader = new FileReader();

   file_element.style.borderRadius = "6px";
   file_element.style.height = "200px";
   file_element.style.margin = "10px";

   function file_reader() {
      reader.onload = () => {
         file_element.setAttribute("src", reader.result);
      }

      reader.readAsDataURL(file);

      file_element.setAttribute("id", `${file_input.files[i].name}`);

      display_div.append(file_element);

      file_element.addEventListener("click", () => {
         media_dialog(file_input.files[i], file_element.getAttribute("src"), file.name, file.type, file.size, i);
      });
   }

   (files_length > 3) ? file_element.style.flex = "33%" : null;

   file_reader();
}

file_input.addEventListener("change", () => {
   if (display_div.hasChildNodes()) {
      let first_el = display_div.firstElementChild;
      while (first_el) {
         first_el.remove();
         first_el = display_div.firstElementChild;
      }
   }

   display_div.style.display = "flex";

   for (let i = 0; i < file_input.files.length; i++) {
      let file_type = file_input.files[i].type,
         file;

      if (file_type.includes("image")) {
         file = document.createElement("img");
         file_format(file, file_input.files[i], file_input.files.length, i);
      } else {
         file = document.createElement("video");
         file_format(file, file_input.files[i], file_input.files.length, i);
      }
   }

   file_input.files.length > 1 || file_input.files.length < 1 ?
      upload_num.innerHTML = `<p>${file_input.files.length} files selected</p>` :
      upload_num.innerHTML = `<p>${file_input.files.length} file selected</p>`;

});

function upload_changer(type, upload_single_color, upload_multiple_color, message) {
   upload_type = type;
   upload_single.style.border = `1px solid ${upload_single_color}`;
   upload_multiple.style.border = `1px solid ${upload_multiple_color}`;
   plain_dialog(message);

   console.log(upload_type);
}

upload_single.addEventListener("click", () => {
   upload_changer("single", "#0096bfab", "grey", "Posts will be shown separately");
});

upload_multiple.addEventListener("click", () => {
   upload_changer("multiple", "grey", "#0096bfab", "All files will be displayed as a single post");
});

upload_btn.addEventListener("click", () => {
   button_loader(upload_btn);

   if (file_input.files.length > 0) {
      let dialog = loading_dialog(),
         upload_res,
         upload_config_res,
         file,
         file_type,
         file_input_size = file_input.files.length,
         default_uploader = (upload_url, file, filename, config_url, pinned, tags, credits, original_name, i) => {
            upload_res = upload(upload_url, file, filename);

            upload_res.then(res => {
               if (res.data.Name !== "") {
                  upload_config_res = upload_config(res.data.Name, config_url, pinned, tags, credits, original_name, upload_type, i);
                  upload_config_res.then(res => {
                     if (res.data.Name === "Upload complete") {
                        file_input_size--;
                        console.log(file_input_size);

                        if (file_input_size == 0) {
                           console.log("DONE");
                           get_req("/generate_new_id").then(res => {
                              if (res.data.Name === "GENERATED") {
                                 dialog.remove();
                                 window.location.reload();
                              }
                           });
                        }
                     }
                  });
               }
            });
         };

      for (let i = 0; i < file_input.files.length; i++) {
         file = file_input.files[i];
         file_type = file.type;

         if (file_type.includes("image")) {
            default_uploader("/upload_meemz", file, "meemz_upload", "/update_meemz_config", upload_caption.value.trim(), upload_tags.value.trim().toLocaleLowerCase(), upload_credits.value.trim(), file.name, i);
         } else {
            default_uploader("/upload_veemz", file, "veemz_upload", "/update_veemz_config", upload_caption.value.trim(), upload_tags.value.trim().toLocaleLowerCase(), upload_credits.value.trim(), file.name, i);
         }
      }

   } else {
      error_dialog("Select files to proceed");
   }
});
