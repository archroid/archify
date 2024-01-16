function upload() {

     // Get the file input element
     var fileInput = document.getElementById('file');
     var progressBar = document.getElementById('uploadProgress');

     for (let i = 0; i < fileInput.files.length; i++) {
          var formData = new FormData();
          formData.append('myFile', fileInput.files[i]);

          var xhr = new XMLHttpRequest();

          // Update progress bar
          xhr.upload.onprogress = function (e) {
               if (e.lengthComputable) {
                    progressBar.max = e.total;
                    progressBar.value = e.loaded;
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
document.getElementById('uploadForm').addEventListener('submit', function (e) {
     e.preventDefault();
     upload();
});


function itemselect(target) {
     const fileNames = Array.from(target.files).map(file => file.name);
     const allNames = fileNames.join("<br>");

     textfile = document.getElementById("file-name");
     textfile.innerHTML = allNames;
     textfile.style.display = "block";


}