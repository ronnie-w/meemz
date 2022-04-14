var qs = Qs;

const search_input_div = document.getElementById("search_input_div");
const search_input = document.getElementById("search_input");
const spinner = document.getElementById("spinner");

let content;
let tags;
let accounts;

function Focus() {
    search_input_div.style.border = "1px dashed white";
}


function FetchRequests() {
    let meemz_eof = document.getElementById("meemz_eof");
    let acc_eof = document.getElementById("acc_eof");
    let tags_eof = document.getElementById("tags_eof");

    meemz_eof.style.display = "none";
    tags_eof.style.display = "none";

    $(".meemz_results_div").empty();
    $(".accounts_results").empty();
    $(".tags_results_div").empty();

    if (search_input.value.trim() !== "") {
        spinner.style.display = "block";

        //search meemz
        axios.post('/search_meemz', qs.stringify({
            ImgOcr: search_input.value.trim(),

        }), {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }).then((res) => {
            if (res.data.length === 0) {
                meemz_eof.innerText = "No Meemz found. Try something else :)";
                meemz_eof.style.display = "grid";
            } else {
                content = res.data;

                content.forEach(c => {
                    $(".meemz_results_div").prepend(
                        `<div style="margin-bottom : 50px ; border-right : 1px solid grey ; border-left : 1px solid grey;" id='meemz_content_main_div ${c.ImgName}'>
                        <div class="uploader_details">
                            <div onclick="Redirect('${c.Username}')">
                                <img class="uploader_profile" alt="uploader_profile" src="/static/profile-pictures/${c.ProfileImg}" loading="lazy" />
                                <p class="uploader_username">${c.Username}</p>
                            </div>
                        </div>
                        <div class="loaded_img_div ${c.ImgName}" style="height : 500px ; filter : blur(5px) ; background-color : #121212;">
                        </div>
                        <img onload="FixHeight('${c.ImgName}')" alt="meemz" class='image ${c.ImgName}' id='image ${c.ImgName}' style="opacity : 0;" loading="lazy"/>
                        <div class="reaction_icons_div">
                            <i class='${c.Reaction1} ${c.ImgName} reaction' onclick="ReactionOnclick('fa-grin-tears' , '${c.ImgName}')" ></i>
                            <i class='${c.Reaction2} ${c.ImgName} reaction' onclick="ReactionOnclick('fa-grin-tongue-squint' , '${c.ImgName}')" ></i>
                            <i class='${c.Reaction3} ${c.ImgName} reaction' onclick="ReactionOnclick('fa-meh' , '${c.ImgName}')" ></i>
                            <i class='${c.Reaction4} ${c.ImgName} reaction' onclick="ReactionOnclick('fa-sad-tear' , '${c.ImgName}')" ></i>
                            <i class='${c.Reaction5} ${c.ImgName} reaction' onclick="ReactionOnclick('fa-angry' , '${c.ImgName}')" ></i>
                        </div>
                        <center>
                        <div class="actions_icons_div">
                            <a href="/static/uploads/${c.ImgName}" download><i class="far fa-arrow-to-bottom"></i></a>
                            <i class="${c.ImgName} far fa-share-alt share" onclick="javascript:PopupToggle('share_popup' , '${c.ImgName}' , 'fa-share-alt' , 'share')" ></i>
                            <i class="${c.ImgName} far fa-comment comment" onclick="javascript:PopupToggle('comment_popup' , '${c.ImgName}' , 'fa-comment' , 'comment'); FetchComments('${c.ImgName}' , '${c.PComment}');" ></i>
                            <i class="${c.ImgName} far fa-flag-alt report" onclick="javascript:PopupToggle('report_popup' , '${c.ImgName}' , 'fa-flag-alt' , 'report')"></i>
                
                        </div>
                        </center>
                
                        <center>
                        <div class='action_popups report_popup ${c.ImgName}' style="display : none ; background-color : rgb(141, 0, 0) ; margin-top : 10px ; width : 90% ; border-radius : 10px ; box-shadow : .5px .5px 7px rgb(141, 0, 0);">
                            <small style="color : white ; background-color : rgb(0 , 0, 0, 0);">Why do you want to report this post ?</small>
                            <select class="report_selector ${c.ImgName}" name="access" id="access" style="width : 90% ; color : red;">
                                <option value="Spam">Report spam</option>
                                <option value="Seen before">Seen this more than once</option>
                                <option value="Not interested">Not interested</option>
                                <option value="Bad quality">Bad quality image or content</option>
                                <option value="Plagiarism">Plagiarism</option>
                                <option value="Adult">Adult content</option>
                            </select>
                            
                            <button style="width : 115px ; background-color : #121212 ; color : red ; font-size : 12px;" onclick="Report('${c.ImgName}')" >Report</button>
                        </div>
                        </center>
                
                        <div style="display : flex ; justify-content : center ; align-items : center ; margin-top : 10px;">
                            <div class='action_popups comment_popup ${c.ImgName}' style="display : none ; background-color : #ffd93b ; margin-top : 10px ; width : 90% ; border-radius : 10px ; box-shadow : .5px .5px 7px #ffd93b;">
                            <div class='comment_popup_pin ${c.ImgName}' style="background-color : transparent ; margin-top : 5px ; width : 100% ; display : flex ; justify-content : center ; align-content : center;"></div>
                            <div style="background-color : transparent ; margin-top : 5px ; width : 100% ; max-height : 200px ; overflow-y : scroll;">
                                <div class='comment_mycomments ${c.ImgName}' style="background-color : transparent ; width : 100% ; display : flex ; flex-direction : column;"></div>
                                <div class='comment_ocomments ${c.ImgName}' style="background-color : transparent ; width : 100% ; display : flex ; flex-direction : column;"></div>
                            </div>
                            <div class='comment_input_div'>
                                <input type="text" class='comment_input ${c.ImgName}' placeholder="Comment..." />
                                <div class='comment_send ${c.ImgName}' onclick="CommentPost('${c.ImgName}' , '${c.PComment}')">
                                    <i class='fal fa-paper-plane ${c.ImgName}' style="color : #121212 ; margin : 10px ; font-size : 15px;"></i>
                                </div>
                            </div>
                            </div>
                        </div>
                
                        <center>
                            <div class='action_popups share_popup ${c.ImgName}' style="display : none ; background-color : white ; margin-top : 10px ; width : 90% ; border-radius : 10px ; box-shadow : .5px .5px 7px white;">
                                <small style="color : #121212 ; background-color : transparent ; font-family: 'Maven Pro', sans-serif;">Share on...</small>
                                <div class="share_icons" style="background-color : transparent ; width : 100% ; display : flex ; justify-content : space-around ; margin-top : 10px ; margin-bottom : 10px;">
                                    <a href='https://twitter.com/intent/tweet?url=https://meemz.gq/public_stats/${c.ImgName}&text=Memes I found on the internet&hashtags=Meemz' style="background-color : transparent;"><i class="fab fa-twitter" style="color : #121212;"></i></a>
                                    <a href='https://reddit.com/submit?url=https://meemz.gq/public_stats/${c.ImgName}&title=Memes I found on the internet' style="background-color : transparent;"><i class="fab fa-reddit" style="color : #121212;"></i></a>
                                    <a href='https://api.whatsapp.com/send?text=Memes I found on the internet%20https://meemz.gq/public_stats/${c.ImgName}' style="background-color : transparent;"><i class="fab fa-whatsapp" style="color : #121212;"></i></a>
                                </div>
                            </div>
                        </center>
                
                    </div>`
                    );

                    let img = document.getElementsByClassName(`image ${c.ImgName}`)[0];
                    let img_div = document.getElementsByClassName(`loaded_img_div ${c.ImgName}`)[0];
                    let observer = new IntersectionObserver(entries => {
                        entries.forEach(entry => {
                            if (entry.isIntersecting === true) {
                                console.log(entry);
                                axios.post("/viewed", qs.stringify({
                                    ImgName: c.ImgName
                                }), {
                                    headers: {
                                        'Content-Type': 'application/x-www-form-urlencoded'
                                    }
                                });

                                img_div.style.filter = "blur(0)";
                                img_div.style.transition = "filter .5s";
                                img_div.style.setProperty("background-color", "none");
                                img.style.opacity = "1";
                                img.setAttribute("src", `/static/uploads/${c.ImgName}`);
                            }
                        });
                    }, {
                        threshold: 0.5
                    });
                    observer.observe(img_div);
                });
            }
            spinner.style.display = "none";
        }).catch((err) => {
            console.log(err)
        });


        //search users
        axios.post('/search_users', qs.stringify({
            Username: search_input.value.trim(),

        }), {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }).then((res) => {
            if (res.data.length === 0) {
                $(".accounts_results").append(`<p id="acc_eof" style="padding-top : 40px ; width : 100% ; display : grid ; place-content : center ; color : grey;">No accounts found. Try something else :)</p>`);
            } else {
                accounts = res.data;

                accounts.forEach(acc => {
                    $(".accounts_results").append(
                        `<div style="border : 1px solid rgb(48, 47, 47) ; margin-top : 10px;">
                        <img onclick="Redirect('${acc.Username}')" src="/static/profile-pictures/${acc.ProfileImg}"
                        alt="convo_banner"
                        style="width : 100px ; height : 100px ; border-radius : 10px ; margin : 5px ; border : 1px solid rgb(48, 47, 47);"
                        />
                        <center>
                            <p style="margin-bottom : 5px;" onclick="Redirect('${acc.Username}')">${acc.Username}
                            </p>
                            <button id="sub ${acc.Username}" style="background-color: ${acc.BGC}; color: ${acc.Color}; font-size: 8px; margin: 2px; box-shadow: .5px .5px 7px wheat;" onclick="Subscription('sub ${acc.Username}' , '${acc.Username}')">${acc.Subscription}</button>
                        </center>
                    </div>`
                    );
                });

            }
        }).catch((err) => {
            console.log(err)
        });
    } else {
        spinner.style.display = "none";
    }

    if (search_input.value.trim() !== "" && search_input.value.length > 2) {
        //search tags
        axios.post('/search_tags', qs.stringify({
            Tags: search_input.value.trim().toLowerCase(),

        }), {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }).then((res) => {
            if (res.data.length === 0) {
                tags_eof.innerText = "No Meemz have this tag. Try something else :)";
                tags_eof.style.display = "grid";
            } else {
                tags = res.data;

                tags.forEach(c => {
                    $(".tags_results_div").prepend(
                        `<div style="margin-bottom : 50px ; border-right : 1px solid grey ; border-left : 1px solid grey;" id='meemz_content_main_div ${c.ImgName}'>
                        <div class="uploader_details">
                            <div onclick="Redirect('${c.Username}')">
                                <img class="uploader_profile" alt="uploader_profile" src="/static/profile-pictures/${c.ProfileImg}" loading="lazy" />
                                <p class="uploader_username">${c.Username}</p>
                            </div>
                        </div>
                        <div class="loaded_img_div ${c.ImgName}" style="height : 500px ; filter : blur(5px) ; background-color : #121212;">
                        </div>
                        <img onload="FixHeight('${c.ImgName}')" alt="meemz" class='image ${c.ImgName}' id='image ${c.ImgName}' style="opacity : 0;" loading="lazy"/>
                        <div class="reaction_icons_div">
                            <i class='${c.Reaction1} ${c.ImgName} reaction' onclick="ReactionOnclick('fa-grin-tears' , '${c.ImgName}')" ></i>
                            <i class='${c.Reaction2} ${c.ImgName} reaction' onclick="ReactionOnclick('fa-grin-tongue-squint' , '${c.ImgName}')" ></i>
                            <i class='${c.Reaction3} ${c.ImgName} reaction' onclick="ReactionOnclick('fa-meh' , '${c.ImgName}')" ></i>
                            <i class='${c.Reaction4} ${c.ImgName} reaction' onclick="ReactionOnclick('fa-sad-tear' , '${c.ImgName}')" ></i>
                            <i class='${c.Reaction5} ${c.ImgName} reaction' onclick="ReactionOnclick('fa-angry' , '${c.ImgName}')" ></i>
                        </div>
                        <center>
                        <div class="actions_icons_div">
                            <a href="/static/uploads/${c.ImgName}" download><i class="far fa-arrow-to-bottom"></i></a>
                            <i class="${c.ImgName} far fa-share-alt share" onclick="javascript:PopupToggle('share_popup' , '${c.ImgName}' , 'fa-share-alt' , 'share')" ></i>
                            <i class="${c.ImgName} far fa-comment comment" onclick="javascript:PopupToggle('comment_popup' , '${c.ImgName}' , 'fa-comment' , 'comment'); FetchComments('${c.ImgName}' , '${c.PComment}');" ></i>
                            <i class="${c.ImgName} far fa-flag-alt report" onclick="javascript:PopupToggle('report_popup' , '${c.ImgName}' , 'fa-flag-alt' , 'report')"></i>
                
                        </div>
                        </center>
                
                        <center>
                        <div class='action_popups report_popup ${c.ImgName}' style="display : none ; background-color : rgb(141, 0, 0) ; margin-top : 10px ; width : 90% ; border-radius : 10px ; box-shadow : .5px .5px 7px rgb(141, 0, 0);">
                            <small style="color : white ; background-color : rgb(0 , 0, 0, 0);">Why do you want to report this post ?</small>
                            <select class="report_selector ${c.ImgName}" name="access" id="access" style="width : 90% ; color : red;">
                                <option value="Spam">Report spam</option>
                                <option value="Seen before">Seen this more than once</option>
                                <option value="Not interested">Not interested</option>
                                <option value="Bad quality">Bad quality image or content</option>
                                <option value="Plagiarism">Plagiarism</option>
                                <option value="Adult">Adult content</option>
                            </select>
                            
                            <button style="width : 115px ; background-color : #121212 ; color : red ; font-size : 12px;" onclick="Report('${c.ImgName}')" >Report</button>
                        </div>
                        </center>
                
                        <div style="display : flex ; justify-content : center ; align-items : center ; margin-top : 10px;">
                            <div class='action_popups comment_popup ${c.ImgName}' style="display : none ; background-color : #ffd93b ; margin-top : 10px ; width : 90% ; border-radius : 10px ; box-shadow : .5px .5px 7px #ffd93b;">
                            <div class='comment_popup_pin ${c.ImgName}' style="background-color : transparent ; margin-top : 5px ; width : 100% ; display : flex ; justify-content : center ; align-content : center;"></div>
                            <div style="background-color : transparent ; margin-top : 5px ; width : 100% ; max-height : 200px ; overflow-y : scroll;">
                                <div class='comment_mycomments ${c.ImgName}' style="background-color : transparent ; width : 100% ; display : flex ; flex-direction : column;"></div>
                                <div class='comment_ocomments ${c.ImgName}' style="background-color : transparent ; width : 100% ; display : flex ; flex-direction : column;"></div>
                            </div>
                            <div class='comment_input_div'>
                                <input type="text" class='comment_input ${c.ImgName}' placeholder="Comment..." />
                                <div class='comment_send ${c.ImgName}' onclick="CommentPost('${c.ImgName}' , '${c.PComment}')">
                                    <i class='fal fa-paper-plane ${c.ImgName}' style="color : #121212 ; margin : 10px ; font-size : 15px;"></i>
                                </div>
                            </div>
                            </div>
                        </div>
                
                        <center>
                            <div class='action_popups share_popup ${c.ImgName}' style="display : none ; background-color : white ; margin-top : 10px ; width : 90% ; border-radius : 10px ; box-shadow : .5px .5px 7px white;">
                                <small style="color : #121212 ; background-color : transparent ; font-family: 'Maven Pro', sans-serif;">Share on...</small>
                                <div class="share_icons" style="background-color : transparent ; width : 100% ; display : flex ; justify-content : space-around ; margin-top : 10px ; margin-bottom : 10px;">
                                    <a href='https://twitter.com/intent/tweet?url=https://meemz.gq/public_stats/${c.ImgName}&text=Memes I found on the internet&hashtags=Meemz' style="background-color : transparent;"><i class="fab fa-twitter" style="color : #121212;"></i></a>
                                    <a href='https://reddit.com/submit?url=https://meemz.gq/public_stats/${c.ImgName}&title=Memes I found on the internet' style="background-color : transparent;"><i class="fab fa-reddit" style="color : #121212;"></i></a>
                                    <a href='https://api.whatsapp.com/send?text=Memes I found on the internet%20https://meemz.gq/public_stats/${c.ImgName}' style="background-color : transparent;"><i class="fab fa-whatsapp" style="color : #121212;"></i></a>
                                </div>
                            </div>
                        </center>
                
                    </div>`
                    );

                    let img = document.getElementsByClassName(`image ${c.ImgName}`)[0];
                    let img_div = document.getElementsByClassName(`loaded_img_div ${c.ImgName}`)[0];
                    let observer = new IntersectionObserver(entries => {
                        entries.forEach(entry => {
                            if (entry.isIntersecting === true) {
                                console.log(entry);
                                axios.post("/viewed", qs.stringify({
                                    ImgName: c.ImgName
                                }), {
                                    headers: {
                                        'Content-Type': 'application/x-www-form-urlencoded'
                                    }
                                });

                                img_div.style.filter = "blur(0)";
                                img_div.style.transition = "filter .5s";
                                img_div.style.setProperty("background-color", "none");
                                img.style.opacity = "1";
                                img.setAttribute("src", `/static/uploads/${c.ImgName}`);
                            }
                        });
                    }, {
                        threshold: 0.5
                    });
                    observer.observe(img_div);
                });
            }
        }).catch((err) => {
            console.log(err)
        });
    } else {
        tags_eof.innerText = "No Meemz have this tag. Try something else :)";
        tags_eof.style.display = "grid";
        tags_eof.style.marginTop = "10px";
        tags_eof.style.marginBottom = "10px";
    }
}

function ReactionOnclick(className, imgName) {
    let el = document.getElementsByClassName(`${imgName} reaction`);
    let main_el = document.getElementsByClassName(`${className} ${imgName} reaction`)[0];
    let main_class = main_el.getAttribute("class");
    if (main_class.substr(0, 3) === "fas") {
        for (let i = 0; i < el.length; i++) {
            let c = el[i].getAttribute("class");
            el[i].setAttribute("class", "far " + c.substr(4));
        }
        axios.post("/delete_reaction", qs.stringify({
            ImgName: imgName
        }), {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }).catch(err => {
            if (err) throw err;
        });
    } else {
        for (let i = 0; i < el.length; i++) {
            let c = el[i].getAttribute("class");
            el[i].setAttribute("class", "far " + c.substr(4));
        }
        main_el.setAttribute("class", `fas ${className} ${imgName} reaction`)

        axios.post("/delete_reaction", qs.stringify({
            ImgName: imgName,

        }), {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }).then(_res => {
            axios.post("/post_reaction", qs.stringify({
                ImgName: imgName,
                ReactionType: className,

            }), {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                }
            }).catch(err => {
                if (err) throw err;
            });
        }).catch(err => {
            if (err) throw err;
        });
    }
};


function ResetToDefaults(imgName) {
    let report = document.getElementsByClassName(`report ${imgName}`)[0];
    let comment = document.getElementsByClassName(`comment ${imgName}`)[0];
    let share = document.getElementsByClassName(`share ${imgName}`)[0];

    report.setAttribute("class", `${imgName} far fa-flag-alt report`);
    comment.setAttribute("class", `${imgName} far fa-comment comment`);
    share.setAttribute("class", `${imgName} far fa-share-alt share`);

    let el = document.getElementsByClassName(`action_popups ${imgName}`);
    for (let i = 0; i < el.length; i++) {
        el[i].style.animation = "fadeOutDown";
        el[i].style.animationDuration = "500ms";
        setTimeout(() => {
            el[i].style.display = "none";
        }, 499);
    }
};

function FetchComments(imgName, pinned_comment) {
    let com_block = document.getElementsByClassName(`comment_popup_pin ${imgName}`)[0];
    let mycomment_div = document.getElementsByClassName(`comment_mycomments ${imgName}`)[0];
    let ocomments_div = document.getElementsByClassName(`comment_ocomments ${imgName}`)[0];
    $(com_block).empty();
    $(mycomment_div).empty();
    $(ocomments_div).empty();
    if (pinned_comment !== "") {
        $(com_block).append(`<i class="fal fa-thumbtack" style="margin-right:10px; font-size:15px; color:#121212;"></i><small style="background-color:transparent;color : #121212;font-family: 'Maven Pro', sans-serif;">${pinned_comment}</small>`)
    } else {
        com_block.style.display = "none";
    }

    axios.post("/fetch_my_comments", qs.stringify({
        ImgName: imgName,

    }), {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }).then(res => {
        let data = res.data;
        if (data[0].Comment !== "") {
            for (let i = 0; i < data.length; i++) {
                $(mycomment_div).prepend(`
                <div style="background-color : transparent;margin-bottom : 10px;" id=${data[i].Comment}>
                    <img src="/static/profile-pictures/${data[i].ProfileImg}" style="
                        width : 35px;
                        height : 35px;
                        border-radius : 200px;
                        float : left;
                        margin-left : 5px;
                        margin-top : 7px;
                        border : 2px solid white;
                    ">
                    <p style="
                        color : grey;
                        float : left;
                        background-color : transparent;
                        font-size : 12px;
                        margin-top : 7px;
                        margin-left : 5px;
                    ">
                        ${data[i].Username}
                    </p>
                    <br/>
                    <div style="
                        width : 100%;
                        background-color : transparent;
                    ">
                        <p style="
                            background-color : transparent;
                            color : black;
                            font-size : 15px;
                            max-width : 300px;
                            margin-top : 5px;
                            word-wrap : break-word;
                            margin-left : 48px;
                        ">
                            ${data[i].Comment}
                        </p>

                    </div>
                </div>
                `);
            }
        }
    });

    axios.post("/fetch_o_comments", qs.stringify({
        ImgName: imgName,

    }), {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
    }).then(res => {
        let data = res.data;
        for (let i = 0; i < data.length; i++) {
            $(ocomments_div).prepend(`
            <div style="background-color : transparent;margin-bottom : 10px;">
                <img src="/static/profile-pictures/${data[i].ProfileImg}" style="
                    width : 35px;
                    height : 35px;
                    border-radius : 200px;
                    float : left;
                    margin-left : 5px;
                    margin-top : 7px;
                    border : 2px solid white;
                ">
                <p style="
                    color : grey;
                    float : left;
                    background-color : transparent;
                    font-size : 12px;
                    margin-top : 7px;
                    margin-left : 5px;
                ">
                    ${data[i].Username}
                </p><br/>
                <div style="
                    width : 100%;
                    background-color : transparent;
                ">
                    <p style="
                        background-color : transparent;
                        color : black;
                        font-size : 15px;
                        max-width : 300px;
                        margin-top : 5px;
                        word-wrap : break-word;
                        margin-left : 48px;
                    ">
                        ${data[i].Comment}
                    </p>

                </div>
            </div>
            `);
        }
    });
};

function CommentPost(imgName, pinned_comment) {
    let comment = document.getElementsByClassName(`comment_input ${imgName}`)[0].value;
    if (comment.length > 0) {
        axios.post("/post_comment", qs.stringify({
            ImgName: imgName,
            ReactionType: comment,

        }), {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        }).then(_res => {
            document.getElementsByClassName(`comment_input ${imgName}`)[0].value = "";
            FetchComments(imgName, pinned_comment)
        });
    }
};

function PopupToggle(popup, imgName, className, activity) {
    let el = document.getElementsByClassName(`${popup} ${imgName}`)[0];
    let clicked_el = document.getElementsByClassName(`${imgName} ${activity}`)[0];
    if (clicked_el.getAttribute("class") === `${imgName} far ${className} ${activity}`) {
        ResetToDefaults(imgName);
        el.style.animation = "fadeInUp";
        el.style.animationDuration = "500ms";
        setTimeout(() => {
            el.style.display = "block";
        }, 500);
        clicked_el.setAttribute("class", `${imgName} far fa-times-circle ${activity}`);
    } else if (clicked_el.getAttribute("class") === `${imgName} far fa-times-circle ${activity}`) {
        ResetToDefaults(imgName);
    }
};

function Redirect(username) {
    axios.post("/fetch_user").then(res => {
        if (username !== res.data.Username) {
            window.location.assign(`/user/${username}`);
        } else {
            window.location.assign("/profile");
        }
    });
};

function FixHeight(image_name) {
    let img_div = document.getElementsByClassName(`loaded_img_div ${image_name}`)[0];
    img_div.style.removeProperty("height");
};

function ToggleSearch(result_div, c_txt) {
    let meemz_results = document.getElementsByClassName("meemz_results")[0];
    let accounts_results = document.getElementsByClassName("accounts_results")[0];
    let tags_results = document.getElementsByClassName("tags_results")[0];

    meemz_results.style.display = "none";
    accounts_results.style.display = "none";
    tags_results.style.display = "none";

    let meemz_c = document.getElementById("meemz_c");
    let tags_c = document.getElementById("tags_c");
    let acc_c = document.getElementById("acc_c");

    meemz_c.style.color = "white";
    tags_c.style.color = "white";
    acc_c.style.color = "white";

    let div = document.getElementsByClassName(result_div)[0];
    let txt = document.getElementById(c_txt);

    txt.style.color = "#121212";
    if (result_div === "accounts_results") {
        div.style.display = "flex";
    } else {
        div.style.display = "block";
    }
}

function Subscription(btn_id, username) {
    let btn = document.getElementById(btn_id);
    if (btn.innerText === "Subscribe") {
        axios.post("/subscribe", qs.stringify({
            Username: username,

        }), {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        });
        btn.style.backgroundColor = "#121212";
        btn.style.color = "wheat";
        btn.innerText = "Unsubscribe";
    } else {
        axios.post("/unsubscribe", qs.stringify({
            Username: username,

        }), {
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            }
        });
        btn.style.backgroundColor = "white";
        btn.style.color = "#121212";
        btn.innerText = "Subscribe";
    }
}