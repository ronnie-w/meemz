import { loading_dialog } from "../modules/dialog.js";

document.querySelector("body").innerHTML +=
    `<div class="meemz_bottom_nav_div">
    <div class="meemz_bottom_nav">
        <i id="meemz_nav_home" class="fi fi-rr-home home"></i>
        <i id="meemz_nav_search" class="fi fi-rr-search search"></i>
        <i id="meemz_nav_add" class="fi fi-rr-add add"></i>
        <i id="meemz_nav_bell" class="fi fi-rr-bell bell"></i>
        <i id="meemz_nav_user" class="fi fi-rr-user user"></i>
    </div>
</div>`;

let home = document.getElementById("meemz_nav_home"),
    search = document.getElementById("meemz_nav_search"),
    create = document.getElementById("meemz_nav_add"),
    notifications = document.getElementById("meemz_nav_bell"),
    profile = document.getElementById("meemz_nav_user");


function click(btn, url) {
    btn.addEventListener("click", () => {
        loading_dialog()
        btn == home ?
            window.location.assign("/")
            : window.location.assign(`/${url}`);
    });
}

click(home), click(search, "search"), click(create, "create"), click(notifications, "notifications"), click(profile, "profile");

var path = window.location.pathname;

switch (path) {
    case "/":
        document.getElementsByClassName("home")[0].setAttribute("class", "fi fi-br-home home");
        break;

    case "/search":
        document.getElementsByClassName("search")[0].setAttribute("class", "fi fi-br-search search");
        break;

    case "/create":
        document.getElementsByClassName("add")[0].setAttribute("class", "fi fi-br-add add");
        break;

    case "/notifications":
        document.getElementsByClassName("bell")[0].setAttribute("class", "fi fi-br-bell bell");
        break;

    default:
        document.getElementsByClassName("user")[0].setAttribute("class", "fi fi-br-user user");
        break;
}