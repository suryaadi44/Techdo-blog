const pathArray = window.location.pathname.split('/');
const id = pathArray[2];

pickBtn = $(".pick-btn");

pickBtn.on("click", async (e) => {
    try {
        let res = await fetch("/post/pick", {
            method: 'POST',
            credentials: 'include',
            body: JSON.stringify({
                postID: id,
            }),
        });
        let response = await res.json();

        if (!response.error) {
            Swal.fire({
                title: 'Post Picked',
                icon: "success",
            });
        } else {
            Swal.fire({
                title: 'Error',
                icon: "error",
                text: response.data,
                showConfirmButton: true,
            });
        }
    } catch (error) {
        Swal.fire({
            title: 'Error',
            icon: "error",
            text: error,
            showConfirmButton: true,
        });
    }
});