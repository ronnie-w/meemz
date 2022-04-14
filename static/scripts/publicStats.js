var qs = Qs;

const path = window.location.pathname;
const locale = path.slice(14, path.length);
const link = locale.split("%20").join(" ");

$(".image").attr("src", `/static/uploads/${link}`);

$("#far_grin").attr("class", `far fa-grin-tears ${link} reaction`);
$("#far_squint").attr("class", `far fa-grin-tongue-squint ${link} reaction`);
$("#far_meh").attr("class", `far fa-meh ${link} reaction`);
$("#far_tear").attr("class", `far fa-sad-tear ${link} reaction`);
$("#far_angry").attr("class", `far fa-angry ${link} reaction`);

axios.post("/fetch_stats", qs.stringify({
    ImgName: link
}), {
    headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
    }
}).then(res => {
    $(".stats_counter1").text(res.data.Reaction1);
    $(".stats_counter2").text(res.data.Reaction2);
    $(".stats_counter3").text(res.data.Reaction3);
    $(".stats_counter4").text(res.data.Reaction4);
    $(".stats_counter5").text(res.data.Reaction5);

    $("#stats_comments").text(res.data.Comments);
    $("#stats_views").text(res.data.Views);
});