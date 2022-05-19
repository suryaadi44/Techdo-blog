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
  $(".edit-profile").show();
  $(".account-settings").hide();
});

navBtn.accountSettingNav.on("click", async (e) => {
  e.preventDefault();
  $(".edit-profile").hide();
  $(".account-settings").show();

  const deleteAccBtn = $(".delete-acc-btn");

  deleteAccBtn.on("click", async (e) => {
    e.preventDefault();
    const { value: confirmation } = await Swal.fire({
      title: "Are you sure?",
      icon: "warning",
      html: '<p>This action cannot be undone. This will permanently delete your account, blog-post, and comments. Please type <b>techdoblog/namaakun</b> to confirm</p>',
      input: "text", 
      confirmButtonText: "I understand the consequences. Delete this account",
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


$(".edit-profile").show();
$(".account-settings").hide();