function ping() {
     const startTime = Date.now();
     fetch(window.location.href + "/ping")
          .then(response => response.json())
          .then(data => {
               const endTime = Date.now();
               const duration = endTime - startTime;
               alert(`${JSON.stringify(data.resp)}, Time: ${duration} ms`);
          })
          .catch((error) => {
               console.error('Error:', error);
          });

}

function accessServer() {
     window.location.href = window.location.href + "/dir/";

}

function reboot() {
     fetch(window.location.href + "/reboot")
          .then(response => response.json())
          .then(data =>
               alert(`${JSON.stringify(data.resp)}`)
          )
          .catch((error) => {
               console.error('Error:', error);
          });

}

function shutdown() {
     fetch(window.location.href + "/shutdown")
          .then(response => response.json())
          .then(data => console.log(data))
          .catch((error) => {
               console.error('Error:', error);
          });
}

function sleep() {
     fetch(window.location.href + "/sleep")
          .then(response => response.json())
          .then(data => console.log(data))
          .catch((error) => {
               console.error('Error:', error);
          });
}

