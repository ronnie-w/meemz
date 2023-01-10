import { default as element_constructor, _ } from "./element.js";

const notification_component = (notification) => {
    let meemz_notification_card = element_constructor("div", "meemz_notification_card"),
        meemz_notification_details = element_constructor("div", "meemz_notification_details"),
        meemz_notification_img = element_constructor("div", "meemz_notification_img"),
        meemz_notification_avi = element_constructor("img", "meemz_notification_avi", _, _, `/static/profile-pictures/${notification.ProfileImg}`),
        meemz_notification_username = element_constructor("div", "meemz_notification_username"),
        username_label = element_constructor("label", _, _, `${notification.Username}`),
        meemz_notification_timestamp = element_constructor("div", "meemz_notification_timestamp"),
        timestamp_label = element_constructor("label", _, _, `${notification.ReceiveTime}`),
        notification_div = element_constructor("div"),
        meemz_notification = element_constructor("p", "meemz_notification", _, `${notification.Notify}`),
        meemz_notifications = document.getElementById("meemz_notifications");

    meemz_notification_img.append(meemz_notification_avi);
    meemz_notification_username.append(username_label);
    meemz_notification_timestamp.append(timestamp_label);

    meemz_notification_details.append(meemz_notification_img, meemz_notification_username, meemz_notification_timestamp);

    notification_div.append(meemz_notification);
    meemz_notification_card.append(meemz_notification_details, notification_div);

    meemz_notifications.append(meemz_notification_card);

    return meemz_notification_card;
}

export default notification_component;