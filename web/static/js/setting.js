function makeDeleteAccountView() {
  return `
    <h1 class="fw-bold ms-5 mt-5 text-danger">Delete Account</h1>
    <div class="delete-account-container mx-5 mt-3 mb-lg-5">
    <p class="delete-msg">Once you delete your account, there is no going back. Please be certain.</p>
    <div class="btn btn-danger delete-acc-btn">Delete Account</div>
    </div>

    <template id="my-template">
      </swal-button>
      <swal-button type="cancel">
        Cancel
      </swal-button>
    </template>
    `
}

function makeEditProfileView() {
  return `
    <h1 class="fw-bold ms-5 mt-5">Edit Profiles</h1>
    <div class="container-fluid pt-5 ps-5 pe-5 pb-5 user-profiles-container">
      <div class="user-profiles__img-container">
        <div class="user-profiles__img" style="background-image: url('/static/img/default-user-profile-img.png');"></div>
        <div class="mb-3 mt-4">
          <label for="banner" class="form-label">Profile Image</label>
          <input class="form-control" type="file" id="banner" name="banner" required>
        </div>
      </div>
      <div class="user-profiles__details">
        <div class="form-outline mb-2">
          <label class="form-label" for="email-form">Email address</label>
          <input class="form-control" id="email-form" type="text" value="alitdarmaputra@gmail.com"
            aria-label="readonly input example" readonly>
        </div>
  
        <div class="form-outline mb-2">
          <label class="form-label" for="username-form">Username</label>
          <input class="form-control" id="username-form" type="text" value="alitdarmaputra"
            aria-label="readonly input example" readonly>
        </div>
  
        <div class="form-outline">
          <div class="row">
            <div class="col">
              <label class="form-label" for="first-name-form">First name</label>
              <input type="text" id="first-name-form" class="form-control" placeholder="First name"
                aria-label="First name">
              <span class="first-name-alert alert" style="font-size: .9em;"></span>
            </div>
            <div class="col">
              <label class="form-label" for="last-name-form">Last name</label>
              <input type="text" id="last-name-form" class="form-control" placeholder="Last name" aria-label="Last name">
              <span class="last-name-alert alert" style="font-size: .9em;"></span>
            </div>
          </div>
        </div>
  
        <div class="form-outline mb-2">
          <label class="form-label" for="username-form">Phone Number</label>
          <input class="form-control" id="phone-form" type="text">
        </div>
  
        <div class="form-outline mb-2">
          <div class="row">
            <label class="form-label" for="about-form">About</label>
            <textarea name="" id="about-form cols=" 10" rows="10"></textarea>
          </div>
        </div>
  
        <div class="button-container d-flex flex-row-reverse">
          <button type="submit" class="btn btn-success">Save</button>
          <button type="submit" class="btn btn-danger me-2">Discard</button>
        </div>
      </div>
    </div>
    `
}

function renderContent(content, thirdPartyContent = null) {
  const container = $(".setting-content");
  container.html("");
  container.html(content);
  container.append(thirdPartyContent);
}

const navBtn = {
  editProfileNav: $(".edit-profile-nav"),
  accountSettingNav: $(".account-setting-nav")
};

navBtn.editProfileNav.on("click", () => {
  renderContent(makeEditProfileView());
});

navBtn.accountSettingNav.on("click", async (e) => {
  e.preventDefault();
  renderContent(makeDeleteAccountView());
  
  const deleteAccBtn = $(".delete-acc-btn");

  deleteAccBtn.on("click", async (e) => {
    e.preventDefault();
    const { value: confirmation } = await Swal.fire({
      title: "Are you sure?",
      icon: "warning",
      html: '<p>This action cannot be undone. This will permanently delete your account, blog-post, and comments. Please type <b>techdoblog/namaakun</b> to confirm</p>',
      input: "text",
      confirmButtonText:"I understand the consequences. Delete this account",
      inputValidator: (value) => {
        if (!value) {
          return 'Confirmation rejected';
        }

        if (value != "techdoblog/namaakun") {
          return 'Confirmation rejected';
        }
      }
    });

    if (confirmation === "techdoblog/namaakun") {
      console.log(confirmation);
      // Do delete account
      Swal.fire({
        title: 'Account Deleted',
        icon: "success",
        showConfirmButton: false,
      }).then((result) => {
          window.location.href = "/" //Location after deleted is success;
      })
    }
  });
});

renderContent(makeEditProfileView());

