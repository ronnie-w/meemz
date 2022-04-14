const home = document.getElementById("home"),
      search = document.getElementById("search"),
      create = document.getElementById("create"),
      convo = document.getElementById("convo"),
      profile = document.getElementById("profile");

const path = window.location.pathname;

function Aqua(ico) {
    ico.style.color = "aqua";
    ico.style.textShadow = ".2px .2px 5px aqua";
}

function White(ico) {
    for (let i = 0; i < ico.length; i++) {
        ico[i].style.color = "white";
    }
}

switch (path) {
    case "/":
        Aqua(home);
        White([search, create, convo, profile]);
        break;
    case "/search":
        Aqua(search);
        White([home, create, convo, profile]);
        break;
    case "/create":
        Aqua(create);
        White([home, search, convo, profile]);
        break;
    case "/convo":
        Aqua(convo);
        White([home, search, create, profile]);
        break;
    case "/profile":
        Aqua(profile);
        White([home, search, create, convo]);
        break;
    default:
        break;
}