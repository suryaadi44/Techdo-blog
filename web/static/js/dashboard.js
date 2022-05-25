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

navMenu.blogBtn.on("click", e => {
  restartState();
  wrappers.blog.show();
});

navMenu.commentBtn.on("click", e => {
  restartState();
  wrappers.comment.show();
});

restartState();
wrappers.dashboard.show();