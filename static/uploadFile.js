function upload() {

     // Get the file input element
     var fileInput = document.getElementById('fileInput');

     // Create a new FormData object
     var formData = new FormData();

     // Add the file to the FormData object
     formData.append('myFile', fileInput.files[0]);
     const startTime = Date.now();

     // Send a POST request with the file data
     fetch("/upload", {
          method: 'POST',
          body: formData
     })
          .then(response => response.json())
          .then(data => {
               const endTime = Date.now();
               const duration = endTime - startTime;
               alert(`${JSON.stringify(data.resp)}, Time: ${duration} ms`);
          })
          .catch((error) => {
               alert('Error:', error);
          });
}
document.getElementById('uploadForm').addEventListener('submit', function (e) {
     e.preventDefault();
     upload();
});

