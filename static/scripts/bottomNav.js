const home = document.getElementById("home"),
      search = document.getElementById("search"),
      create = document.getElementById("create"),
      veemz = document.getElementById("veemz"),
      profile = document.getElementById("profile");

const path = window.location.pathname;

function Aqua(ico) {
    $(ico).removeClass("fal");
    $(ico).addClass("far");
}

function White(ico) {
    for (let i = 0; i < ico.length; i++) {
        ico[i].style.color = "#121212";
    }
}

switch (path) {
    case "/":
        Aqua(home);
        White([search, create, veemz, profile]);
        break;
    case "/search":
        Aqua(search);
        White([home, create, veemz, profile]);
        break;
    case "/create":
        Aqua(create);
        White([home, search, veemz, profile]);
        break;
    case "/veemz":
        Aqua(veemz);
        White([home, search, create, profile]);
        break;
    case "/profile":
        Aqua(profile);
        White([home, search, create, veemz]);
        break;
    default:
        break;
}