function upload() {
    // Get the file input element
    var fileInput = document.getElementById('file');
    var progressBar = document.getElementById('uploadProgress');

    // Calculate total size of all files
    var totalSize = 0;
    for (let i = 0; i < fileInput.files.length; i++) {
        totalSize += fileInput.files[i].size;
    }

    // Keep track of total uploaded size
    var totalLoaded = 0;

    for (let i = 0; i < fileInput.files.length; i++) {
        var formData = new FormData();
        formData.append('myFile', fileInput.files[i]);

        var xhr = new XMLHttpRequest();

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
    const allNames = fileNames.join("<br>");

    textfile = document.getElementById("file-name");
    textfile.innerHTML = allNames;
    textfile.style.display = "block";


}