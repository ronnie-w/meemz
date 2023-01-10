let rating_up = document.getElementById("home_rating_up"),
    rating_down = document.getElementById("home_rating_down");
    async function animator (target, animation, animationDuration) {
        target.removeAttribute("style");

        window.setTimeout(() => {
            target.style.animation = animation;
            target.style.animationDuration = animationDuration;
        }, 10);
    };

console.log(rating_up);
rating_up.addEventListener("click", () => {
    animator(rating_up, "fadeInUp", "500ms");
});

rating_down.addEventListener("click", () => {
    animator(rating_down, "fadeInDown", "500ms");
});