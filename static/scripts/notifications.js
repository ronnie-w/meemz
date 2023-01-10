import { get_req } from "./modules/dialer.js";
import notification_component from "./components/notification.js";

let start_index = 0,
    meemz_notifications_count = document.getElementById("meemz_notifications_count");

(_ => {
    get_req("/fetch_notifications").then(res => {
        res.data.length == 0 || res.data == undefined 
            ? meemz_notifications_count.innerText = "You have no new notifications" 
            : res.data.length == 1 
                ? meemz_notifications_count.innerText = `You have ${res.data.length} unread notification` 
                : meemz_notifications_count.innerText = `You have ${res.data.length} unread notifications`;
        res.data = res.data.reverse();
        function DisplayNotification() {
            let notification = notification_component(res.data[start_index]),
                observer = new IntersectionObserver(
                    entries => {
                        entries.map(entry => {
                            if (entry.isIntersecting) {
                                DisplayNotification();
                                observer.unobserve(notification);
                            }
                        });
                    }, {
                    threshold: 1.0
                }
                );

            if (start_index < res.data.length - 1) {
                start_index++;
                observer.observe(notification);
            }
        }

        DisplayNotification();
    });
})();