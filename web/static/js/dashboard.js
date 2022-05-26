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
  return `
    <div class="row card blog-items my-3">
      <div class="row py-3 px-3">
        <div class="col align-self-center">
          <h5>${blog.Title}</h5>
        </div>

        <div class="col mt-2 align-self-center">
          <p class="text-muted">${blog.CreatedAt}</p>
        </div>

        <div class="col-3 align-self-center">
          <div data-blogId = "${blog.PostID}" class="btn btn-warning edit-btn mb-2"><i class="fa-regular fa-pen-to-square"></i></div>
          <div data-blogId = "${blog.PostID}" class="btn btn-danger delete-btn mb-2"><i class="fa-regular fa-trash-can"></i></div>
        </div>
      </div>
    </div>
  `
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
  
  for (const blog of blogData) {
    blogListContainer.append(makeBlogItems(blog))
  }
  wrappers.blog.show();
});

function makeCommentItem(comment) {
  return `
  <div class="row card comment-items my-3">
    <div class="row py-3 px-3">
      <div class="col align-self-center">
        <h5>Blog Title</h5>
      </div>

      <div class="col mt-2 align-self-center">
        <p class="text-muted">${comment.commentBody}</p>
      </div>

      <div class="col-3 align-self-center">
        <div data-commentID="${comment.index}" class="btn btn-danger delete-btn mb-2"><i class="fa-regular fa-trash-can"></i></div>
      </div>
    </div>
  </div>
  `
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
  for(const comment of commentData) {
    commentListContainer.append(makeCommentItem(comment));
  }
  wrappers.comment.show();
});

restartState();
wrappers.dashboard.show();