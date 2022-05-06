$('#summernote').summernote({
    height: ($(window).height() - 400),
    callbacks: {
        onImageUpload: function(image) {
            uploadImage(image[0]);
        },
        onMediaDelete : function(target) {
            deleteImage(target[0].src);
        }
    }
});

function uploadImage(image, editor) {
    var data = new FormData();
    data.append("image", image);
    $.ajax({
        url: '/api/upload/image',
        cache: false,
        contentType: false,
        processData: false,
        data: data,
        type: "post",
        success: function(url) {
            let image = $('<img>').attr('src', url);
            editor.summernote("insertNode", image[0]);
        },
        error: function(data) {
        }
    });
}

function deleteImage(src) {
    $.ajax({
        data: {src : src},
        type: "POST",
        url: "YOUR DELETE SCRIPT", 
        cache: false,
        success: function(data) {
            alert(data);
        }
    });
}