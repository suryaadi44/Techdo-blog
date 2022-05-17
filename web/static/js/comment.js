const comments = {
    commentCount: $(".comment-count"),
    commentWrapper: $(".comment-wrapper"),
    commentAlert: $(".comment-alert"),
    commentForm: $("#comment-form"),
}

function makeCommentCard(data) {
    if(data)
    return `
    <div class="comment-list mb-3 row">
        <div class="col-3 user-pict" style="background-image: url('${data.userPic}');">
        </div>

        <div class="col card mx-2 px-4 py-3 comment-card">
            <h5 class="comment-header">${data.firstName} ${data.lastName} • ${data.createdAt}</h5>
            <p class="comment-body">${data.commentBody}</p>
        </div>
    </div>`
}

function checkForm(commentBody) {
    comments.commentAlert.html("");
    
    if(commentBody)
        return true;
    
    comments.commentAlert.html("Please provide comment");
    return false;
}

function getComment() {
    let commentList;
    let commentItems = "";
    
    // comments.commentWrapper.html(commentItems);

    const route = "/post/" + comments.commentWrapper.data("postid") + "/comment";

    fetch(route)
    .then((response) => response.json())
    .then((result) => {
        let commentCount = 0;
        if(result) {
            commentList = result.data;
            if(commentList) {
                commentCount = commentList.length;
                for(let i = 0; i<commentCount; i++) {
                    commentItems += makeCommentCard(commentList[i]);
                }
                comments.commentWrapper.html(commentItems);
            }    
        }
        comments.commentCount.html(`${commentCount}`);
    });
}

function addComment() {
    if(checkForm(comments.commentForm.val())) {
        const post_id = comments.commentWrapper.data("postid");
        const route = "/post/" + post_id + "/comment/add";
      
        fetch(route, {
            method: "POST",
            headers: {
                Accept: "application/json, text/plain, */*",
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                "commentBody": comments.commentForm.val()
            })
        })
            .then((response) => response.json())
            .then((result) => {
                if(result) {
                    getComment();
                }
            })
    }
}

$(".comment-btn").on("click", () => addComment());
comments.commentWrapper.html("");
getComment();
