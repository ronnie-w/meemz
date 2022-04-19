if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('/sw.js').catch(
      err => {
        if (err) throw err;
      });
} else {
    console.log("Unsupported application");
}