$("#post-form").submit(function (e) {
    e.preventDefault();

    document.getElementById("submit-btn").disabled = true
    document.getElementById("submit-btn").innerHTML = `<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> Post`;

    const fileInput = document.getElementById('banner');
    const formData = new FormData();

    formData.append('banner', fileInput.files[0]);
    formData.append('title', document.getElementById('title').value);
    formData.append('category', document.getElementById('category').value);
    formData.append('editordata', document.getElementById('summernote').value);

    fetch('/post/create', {
        method: 'POST',
        body: formData,
    }).then((response) => response.json())
        .then((result) => {
            if (result.error) {
                Swal.fire({
                    icon: "error",
                    title: "Upload error!",
                    text: result.data,
                    confirmButtonText: "Try again",
                    allowOutsideClick: false,
                    allowEscapeKey: false,
                });
            } else {
                document.getElementById("submit-btn").innerHTML = "Post"
                Swal.fire({
                    title: "Success",
                    text: "Success creating blog post",
                    icon: "success",
                    allowOutsideClick: false,
                    allowEscapeKey: false,
                    confirmButtonText: "Visit page",
                }).then(function () {
                    window.location = result.data;
                });
            }
        });
})