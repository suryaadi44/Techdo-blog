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
    wrapper.css("display", "none");
  }
}

navMenu.dashboardBtn.on("click", e => {
  restartState();
  e.css("display", "block");
});

navMenu.blogBtn.on("click", e => {
  restartState();
  e.css("display", "block");
});

navMenu.commentBtn.on("click", e => {
  restartState();
  e.css("display", "block");
});