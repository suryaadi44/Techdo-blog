const navMenu = {
    dashboardBtn: $(".dashboard-btn"),
    blogBtn: $(".blog-btn"),
    commentBtn: $(".comment-btn")
}

const viewContainer = $(".main-container");
let viewToRender;

function renderView(viewToRender) {
    viewContainer.html("");
    viewContainer.html(viewToRender);
}

function makeDashboardView() {
    return `
    <!-- Dashboard -->
    <h2>Dashboard</h2>
    <div class="row statistic-container">
      <div class="col-md me-3 card post-statistic-container" style="width: 100%;">
        <div class="card-body">
          <h4 class="card-title">Total Post</h4>
          <div class="row py-1 statistic-container">
            <h2 class="total-post col">12</h2>
            <div class="col">
              <i class="fa-regular fa-file-lines"></i>
            </div>
          </div>
          <p class="total-post col">Posted blog</p>
        </div>
      </div>

      <div class="col-md blog-post-list card" style="width: 100%;">
        <div class="card-body">
          <h4 class="card-title">Total Comment</h4>
          <div class="row py-1 statistic-container">
            <h2 class="total-comment col">12</h2>
            <div class="col">
              <i class="fa-regular fa-comment-dots"></i>
            </div>
          </div>
          <p class="total-post col">Posted blog</p>
        </div>
      </div>

    </div>

    <div class="row">
      <div class="mt-3 blog-post-list card" style="width: 100%;">
        <div class="card-body">
          <h4 class="card-title">Recent Post</h4>
          <!-- Post List -->
          <div class="post-list-container">
            <table class="table">
              <thead>
                <tr>
                  <th scope="col">No.</th>
                  <th scope="col">Title</th>
                  <th scope="col">Post Date</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <th scope="row">1</th>
                  <td>Mark</td>
                  <td>Otto</td>
                </tr>
                <tr>
                  <th scope="row">2</th>
                  <td>Jacob</td>
                  <td>Thornton</td>
                </tr>
                <tr>
                  <th scope="row">3</th>
                  <td>Larry</td>
                  <td>the Bird</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <div class="row">
      <div class="mt-3 blog-post-list card" style="width: 100%;">
        <div class="card-body">
          <h4 class="card-title">Recent Comment</h4>
          <!-- Comment List -->
          <div class="post-list-container">
            <table class="table">
              <thead>
                <tr>
                  <th scope="col">No.</th>
                  <th scope="col">Title</th>
                  <th scope="col">Comment</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <th scope="row">1</th>
                  <td>Mark</td>
                  <td>Otto</td>
                </tr>
                <tr>
                  <th scope="row">2</th>
                  <td>Jacob</td>
                  <td>Thornton</td>
                </tr>
                <tr>
                  <th scope="row">3</th>
                  <td>Larry</td>
                  <td>the Bird</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
    `
}

function makeBlogView() {
    return `
    <!-- Your Blog -->

    <h2>Your Blog</h2>
    <div class="col your-blog-list">

      <div class="row card blog-items my-3">
        <div class="row py-3 px-3">
          <div class="col align-self-center">
            <h5>Blog Title</h5>
          </div>

          <div class="col mt-2 align-self-center">
            <p class="text-muted">12 April 2022</p>
          </div>

          <div class="col-3 align-self-center">
            <div class="btn btn-warning edit-btn mb-2"><i class="fa-regular fa-pen-to-square"></i></div>
            <div class="btn btn-danger delete-btn mb-2"><i class="fa-regular fa-trash-can"></i></div>
          </div>
        </div>
      </div>
    </div>
    `
}

function makeCommentView() {
    return `
    <!-- Your Comment -->

    <h2>Your Comment</h2>
    <div class="col your-blog-list">

      <div class="row card blog-items my-3">
        <div class="row py-3 px-3">
          <div class="col align-self-center">
            <h5>Blog Title</h5>
          </div>

          <div class="col mt-2 align-self-center">
            <p class="text-muted">Lorem ipsum dolor sit amet consectetur adipisicing elit. Assumenda, quia hic? Natus eveniet deserunt aliquam in doloremque tempora. Eos ut animi omnis similique voluptatum deleniti, tempora eveniet accusamus ipsa illo?</p>
          </div>

          <div class="col-3 align-self-center">
            <div class="btn btn-danger delete-btn mb-2"><i class="fa-regular fa-trash-can"></i></div>
          </div>
        </div>
      </div>
    </div> 
    `
}

navMenu.dashboardBtn.on("click", e => {
    renderView(makeDashboardView());
});

navMenu.blogBtn.on("click", e => {
    renderView(makeBlogView());
});

navMenu.commentBtn.on("click", e => {
    renderView(makeCommentView());
});

renderView(makeDashboardView());