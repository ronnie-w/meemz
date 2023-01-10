var match = document.cookie.match(new RegExp('(^| )color-scheme=([^;]+)'));

const theme_toggle = (theme) => {
    document.getElementById("color-scheme").href = `/static/stylesheets/themes/${theme}.css`;
    document.cookie = `color-scheme=${theme}; expires=Fri, 18 Jan 2050 12:00:00 UTC; path=/; samesite=lax`;
}

match && !match[2].includes("auto") ?
    match[2].includes("dark") ? theme_toggle("dark") : theme_toggle("light")
    : theme_toggle("auto");

export default theme_toggle;