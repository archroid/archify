const toggler = document.getElementById('theme-toggle');
var metaThemeColor = document.querySelector("meta[name=theme-color]");



function changeThemeColor(isDark) {
     if (isDark) {
          metaThemeColor.setAttribute("content", "#181a1e");
     } else {
          metaThemeColor.setAttribute("content", "#eeeeee");
     }
}


if (localStorage.getItem('isDarkMode') === "false") {
     toggler.checked = false;
     document.body.classList.remove('dark');
     changeThemeColor(false);

} else {
     toggler.checked = true;
     document.body.classList.add('dark');
     changeThemeColor(true);

}

toggler.addEventListener('change', function () {
     if (this.checked) {
          document.body.classList.add('dark');
          localStorage.setItem('isDarkMode', true);
          changeThemeColor(true);

     } else {
          document.body.classList.remove('dark');
          localStorage.setItem('isDarkMode', false);
          changeThemeColor(false);
     }
});



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
     butterup.toast({
          title: 'Rebooting!',
          location: 'top-center',
          icon: true,
          dismissable: true,
          type: 'success',
          // theme: 'glass',
     });

     fetch(window.location.href + "/reboot")
          .then(response => response.json())
          .then(data => console.log(data))
          .catch((error) => {
               // console.error('Error:', error);
               // butterup.toast({
               //      title: 'Error rebooting!',
               //      location: 'top-center',
               //      icon: true,
               //      dismissable: true,
               //      type: 'error',
               // });
          });

}

function shutdown() {

     butterup.toast({
          title: 'Shutting down!',
          location: 'top-center',
          icon: true,
          dismissable: true,
          type: 'success',
          // theme: 'glass',
     });


     fetch(window.location.href + "/shutdown")
          .then(response => response.json())
          .then(data => console.log(data))
          .catch((error) => {
               // butterup.toast({
               //      title: 'Error shutting down!',
               //      location: 'top-center',
               //      icon: true,
               //      dismissable: true,
               //      type: 'error',
               // });
          });
}

function sleep() {

     butterup.toast({
          title: 'Sleeping!',
          location: 'top-center',
          icon: true,
          dismissable: true,
          type: 'success',
          // theme: 'glass',
     });

     fetch(window.location.href + "/sleep")
          .then(response => response.json())
          .then(data => console.log(data))
          .catch((error) => {
               // butterup.toast({
               //      title: 'Error sleeping!',
               //      location: 'top-center',
               //      icon: true,
               //      dismissable: true,
               //      type: 'error',
               // });
          });
}


function upload() {
     // Get the file input element
     var fileInput = document.getElementById('file');
     var progressBar = document.getElementById('uploadProgress');

     // Calculate total size of all files
     var totalSize = 0;
     for (let i = 0; i < fileInput.files.length; i++) {
          totalSize += fileInput.files[i].size;
     }

     if (totalSize == 0) {
          butterup.toast({
               title: 'No files selected!',
               location: 'top-center',
               icon: true,
               dismissable: true,
               type: 'error',
               // theme: 'glass',
          });
     }

     // Keep track of total uploaded size
     var totalLoaded = 0;

     for (let i = 0; i < fileInput.files.length; i++) {
          var formData = new FormData();
          formData.append('myFile', fileInput.files[i]);

          var xhr = new XMLHttpRequest();

          progressBar.style.visibility = "visible";

          // Update progress bar
          xhr.upload.onprogress = function (e) {
               if (e.lengthComputable) {
                    totalLoaded += e.loaded;
                    progressBar.max = totalSize;
                    progressBar.value = totalLoaded;
               }
          };

          // Load end
          xhr.onloadend = function () {
               if (xhr.status == 200) {
                    console.log("upload complete");
               } else {
                    console.log("upload failed");
               }
          };

          xhr.open('POST', '/upload', true);
          xhr.send(formData);
     }
}


function itemselect(target) {
     const fileNames = Array.from(target.files).map(file => file.name);

     // Assuming you have a table with id "myTable"
     var table = document.getElementById("name-table");
     table.innerHTML = '';
     // Loop through the array
     for (var i = 0; i < fileNames.length; i++) {
          // Create a new row
          var row = table.insertRow();

          // Create a new cell for the file type
          var typeCell = row.insertCell();

          // Determine the file type based on the file extension
          var fileType = fileNames[i].split('.').pop().toLowerCase();
          console.log(fileType);

          // Choose an emoji based on the file type
          var emoji;
          switch (fileType) {
               case 'mp4':
               case 'mkv':
               case 'avi':
               case 'mov':
               case 'flv':
               case 'wmv':
               case 'webm':
                    emoji = '🎥';
                    break;
               case 'jpg':
               case 'jpeg':
               case 'gif':
               case 'png':
               case 'bmp':
                    emoji = '📸';
                    break;
               case 'mp3':
               case 'wav':
               case 'ogg':
               case 'flac':
                    emoji = '🎵';
                    break;
               case 'zip':
               case 'rar':
               case '7z':
               case 'tar':
               case 'gz':
               case 'bz2':
               case 'xz':
                    emoji = '📦';
                    break
               default:
                    emoji = '📄';
          }

          // Insert the emoji into the type cell
          typeCell.textContent = emoji;

          // Create a new cell for the filename
          var nameCell = row.insertCell();

          // Insert the filename into the name cell
          nameCell.textContent = fileNames[i];
     }

}

function kill() {
     window.location.href = window.location.href + "/kill"
}

function log() {
     window.location.href = window.location.href + "/log"
}