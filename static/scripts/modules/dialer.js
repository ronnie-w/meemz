'use strict'

const qs = Qs;

const post_form = (url, data) => {
    let res = axios.post(url, qs.stringify(data), {
        headers: {
            "Content-Type": "application/x-www-form-urlencoded"
        }
    });

    return res;
}

const get_req = (url) => {
    let res = axios.get(url)

    return res;
}

const upload = (upload_url, file, filename) => {
    let formData = new FormData(),
        res;

    formData.append(filename, file);

    res = axios.post(upload_url, formData, {
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });

    return res;
}

const upload_config = (file_name, config_url, pinned, tags, credits, original_name, upload_type, image_index) => {
    let res = axios.post(config_url, qs.stringify({
        Credits: credits,
        Pinned: pinned.trim(),
        Tags: tags.trim().toLowerCase(),
        FileName: file_name,
        OriginalName: original_name,
        UploadType: upload_type,
        FileIndex: image_index
    }), {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    });

    return res;
}

export { post_form as default, get_req, upload, upload_config };
