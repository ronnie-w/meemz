self.addEventListener('install', (e) => {
    e.waitUntil(
        caches.open("templates").then(cache => {
            return cache.addAll([
                "/static/stylesheets",
                "/static/scripts",
                "/static/icons",
                "./templates"
            ]);
        })
    )
});

self.addEventListener('fetch', (_e) => {
    return;
});
