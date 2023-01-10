const cache_name = "meemz_cache-v1",
        stylesheets_path = "static/stylesheets",
        urls_to_cache = [
            "static/scripts/dependencies",
            "/", "/search", "/create", "/notifications", "/profile",
            "manifest.json",
            "sw.js",
            "sw_init.js",
            stylesheets_path,
            stylesheets_path + "/themes",
            stylesheets_path + "/fonts",
            stylesheets_path + "/embedded-opentypes",
            stylesheets_path + "/icons",
            stylesheets_path + "/woff",
            stylesheets_path + "/woff2",
            "static/scripts/components",
            "static/scripts/modules",
            "static/scripts",
            "static/icons",
            "static/profile-pictures",
        ]

self.addEventListener('install', (e) => {
    console.log("Installing service worker");
    e.waitUntil(
        caches.open(cache_name).then(cache => {
            return cache.addAll(urls_to_cache);
        })
    )
});

self.addEventListener('activate', e => {
    console.log("Activating service worker", e);
});

self.addEventListener("fetch", e => {
    e.respondWith(caches.match(e.request).then(cache_response => {
        return cache_response || fetch(e.request);
    }));
});

self.addEventListener('push', e => {
    const title = e.data.text();
    e.waitUntil(self.registration.showNotification(title));
});