if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('/sw.js').then(registration => { console.log("SW registration scope: ", registration.scope); }).catch(
    err => {
      if (err) throw err;
    });

  navigator.serviceWorker.ready.then(reg => {
    reg.pushManager.subscribe({ 
      userVisibleOnly: true, 
      applicationServerKey: 'BEl62iUYgUivxIkv69yViEuiBIa-Ib9-SkvMeAtA3LFgDzkrxZJjSgSnfckjBJuBkr3qBUYIHBQFLXYp5Nksh8U'
        }).then(sub => {
      //send sub.Json to server
      console.log(sub.toJSON());
    });
  });
} else {
  console.log("Unsupported application");
}