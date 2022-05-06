$('#summernote').summernote({
    height: ($(window).height() - 400),
    callbacks: {
        onImageUpload: function (image) {
            uploadImage(image[0]);
        },
        onMediaDelete: function (target) {
            deleteImage(target[0].id);
        }
    }
});

function uploadImage(image) {
    var data = new FormData();
    data.append("image", image);

    const uploadUrl = "/api/upload/image"

    fetch(uploadUrl, {
        method: "post",
        body: data,
    })
        .then((response) => response.json())
        .then((result) => {
            if (!result.error) {
                let image = $('<img>').attr('src', result.data.url).attr('id', result.data.fileId).attr("readonly", "readonly");
                $('#summernote').summernote("insertNode", image[0]);
            }
            // TODO: Add indication if upload error
        });
}

function deleteImage(id) {
    const deleteUrl = "/api/delete/image"
    fetch(deleteUrl, {
        method: "post",
        body: JSON.stringify({
            fileId: id,
        }),
    })
        .then((response) => response.json())
        .then((result) => {
            // TODO: Add indication if delete error
        });
}