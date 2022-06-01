const navMenu = {
  dashboardBtn: $(".dashboard-btn"),
  blogBtn: $(".blog-btn"),
  commentBtn: $(".comment-btn")
}

const wrappers = {
  dashboard: $(".dashboard-wrapper"),
  blog: $(".blog-wrapper"),
  comment: $(".comment-wrapper"),
}

function restartState() {
  for (const wrapper in wrappers) {
    wrappers[wrapper].hide();
  }
}

navMenu.dashboardBtn.on("click", e => {
  restartState();
  wrappers.dashboard.show();
});

function makeBlogItems(blog) {
  if (blog.title) {
    return `
      <div class="row card blog-items my-3">
        <div class="row py-3 px-3">
          <div class="col align-self-center">
            <h5>
              <a href="/post/${blog.postId}" class="post-link">${blog.title}</a>
            </h5>
          </div>
  
          <div class="col mt-2 align-self-center">
            <p class="text-muted">${blog.createdAt}</p>
          </div>
  
          <div class="col-3 align-self-center">
            <div data-blogId = "${blog.postId}" class="btn btn-warning edit-post-btn mb-2">
              <i class="fa-regular fa-pen-to-square"></i>
            </div>
            <div data-blogId = "${blog.postId}" class="btn btn-danger delete-post-btn mb-2">
              <i class="fa-regular fa-trash-can"></i>
            </div>
          </div>
        </div>
      </div>
    `
  }

  return `
    <div class="row card blog-items my-3">
      <div class="row py-3 px-3">
        <div class="col align-self-center">
          <h5>
            No Post
          </h5>
        </div>
      </div>
    </div>
  `
}

function makeCategoryItem(category) {
  return `
    <h4 style="text-transform: capitalize;">${category}</h4>
  `
}

async function listenPostEvent() {
  $(".delete-post-btn").each(async function () {
    $(this).on("click", async () => {
      id = $(this).attr("data-blogId");

      await Swal.fire({
        title: "Are you sure?",
        text: "You won't be able to revert this!",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonText: 'Delete',
        confirmButtonColor: '#d33',
      }).then((result) => {
        if (result.isConfirmed) {
          let res = fetch("/post/" + id + "/delete", {
            method: "DELETE",
            credentials: 'include',
          }).then((response) => response.json())
            .then((result) => {
              if (result.error) {
                Swal.fire({
                  icon: "error",
                  title: "Delete error!",
                  text: result.data,
                  confirmButtonText: "Ok",
                  allowOutsideClick: false,
                  allowEscapeKey: false,
                });
              } else {
                Swal.fire({
                  title: "Success",
                  text: "Success deleting blog post",
                  icon: "success",
                  allowOutsideClick: false,
                  allowEscapeKey: false,
                  confirmButtonText: "Ok",
                })
                restartState()
                wrappers.dashboard.show();
              }
            });
        }
      })
    });
  });

  $(".edit-post-btn").each(async function () {
    $(this).on("click", async () => {
      id = $(this).attr("data-blogId");
      window.location.href = "/post/" + id + "/edit"
    });
  });
}

navMenu.blogBtn.on("click", async e => {
  restartState();
  const blogListContainer = $(".your-blog-list");
  blogListContainer.html("");
  // Fetch blog data from backend
  let blogData = await fetch("/user/post", {
    method: "GET",
    headers: {
      'Content-Type': 'application/json'
    },
  }, (response) => JSON.stringify(response))
    .then((result) => result.json());

  blogData = blogData.data;

  let dict = {}
  for (const blog of blogData) {
    if (!blog.category) {
      blog.category = "No Category"
    }

    if (dict[blog.category]) {
      dict[blog.category].push(blog);
    } else {
      dict[blog.category] = [blog]
    }
  }

  for (const [key, value] of Object.entries(dict)) {
    blogListContainer.append(makeCategoryItem(key))
    for (const post of value) {
      blogListContainer.append(makeBlogItems(post))
    }
  }

  wrappers.blog.show();
  listenPostEvent();
});

function makeCommentItem(comment) {
  return `
  <div class="row card comment-items my-3">
    <div class="row py-3 px-3">
      <div class="col align-self-center">
        <h5>
        <a href="/post/${comment.postID}" class="post-link">${comment.postTitle}</a>
        </h5>
      </div>

      <div class="col mt-2 align-self-center">
        <p class="text-muted">${comment.commentBody}</p>
      </div>

      <div class="col-3 align-self-center">
        <div data-commentID="${comment.commentID}" class="btn btn-danger delete-comment-btn mb-2">
          <i class="fa-regular fa-trash-can"></i>
        </div>
      </div>
    </div>
  </div>
  `
}

async function listenCommentEvent() {
  $(".delete-comment-btn").each(async function () {
    $(this).on("click", async () => {
      id = $(this).attr("data-commentID");
      await Swal.fire({
        title: "Are you sure?",
        text: "You won't be able to revert this!",
        icon: 'warning',
        showCancelButton: true,
        confirmButtonText: 'Delete',
        confirmButtonColor: '#d33',
      }).then((result) => {
        if (result.isConfirmed) {
          let res = fetch("/post/comment/delete", {
            method: "DELETE",
            credentials: 'include',
            body: JSON.stringify({
              commentID: id,
            }),
          }).then((response) => response.json())
            .then((result) => {
              if (result.error) {
                Swal.fire({
                  icon: "error",
                  title: "Delete error!",
                  text: result.data,
                  confirmButtonText: "Ok",
                  allowOutsideClick: false,
                  allowEscapeKey: false,
                });
              } else {
                Swal.fire({
                  title: "Success",
                  text: "Success deleting comment",
                  icon: "success",
                  allowOutsideClick: false,
                  allowEscapeKey: false,
                  confirmButtonText: "Ok",
                })
                restartState()
                wrappers.dashboard.show();
              }
            });
        }
      })
    });
  });
}

navMenu.commentBtn.on("click", async e => {
  restartState();
  const commentListContainer = $(".your-comment-list");
  commentListContainer.html("");
  // Fetch comment data from backend
  let commentData = await fetch("/user/comment", {
    method: "GET",
    headers: {
      'Content-Type': 'application/json'
    },
  }, (response) => JSON.stringify(response))
    .then((result) => result.json());

  commentData = commentData.data;
  for (const comment of commentData) {
    commentListContainer.append(makeCommentItem(comment));
  }
  wrappers.comment.show();
  listenCommentEvent();
});

restartState();
wrappers.dashboard.show();