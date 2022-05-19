async function getUserMiniDetail() {
  try{
    let res = await fetch("/user/mini-detail", {
      method: 'GET',
      credentials: 'include'
    });
    return await res.json();
  } catch (error) {
    Swal.fire({
      title: 'Error',
      icon: "error",
      text: error,
      showConfirmButton: true,
    });
  }
}

async function deleteAccount(){
  try {
    let res = await fetch("/user/delete", {
      method: 'DELETE',
      credentials: 'include'
    });
    return await res.json();
  } catch (error) {
    Swal.fire({
      title: 'Error',
      icon: "error",
      text: error,
      showConfirmButton: true,
    });
  }

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
  $(".edit-profile").show();
  $(".account-settings").hide();
});

navBtn.accountSettingNav.on("click", async (e) => {
  e.preventDefault();

  let user = await getUserMiniDetail();

  $(".edit-profile").hide();
  $(".account-settings").show();

  const deleteAccBtn = $(".delete-acc-btn");

  deleteAccBtn.on("click", async (e) => {
    e.preventDefault();
    const { value: confirmation } = await Swal.fire({
      title: "Are you sure?",
      icon: "warning",
      html: '<p>This action cannot be undone. This will permanently delete your account, blog-post, and comments. Please type <b>techdoblog/' + user.data.username + '</b> to confirm</p>',
      input: "text",
      confirmButtonText: "I understand the consequences. Delete this account",
      inputValidator: (value) => {
        if (!value) {
          return 'Confirmation rejected';
        }
        if (value != "techdoblog/" + user.data.username) {
          return 'Confirmation rejected';
        }
      }
    });

    if (confirmation === "techdoblog/" + user.data.username) {
      // Do delete account
      let result = await deleteAccount()
      
      if (!result.error) {
        Swal.fire({
          title: 'Account Deleted',
          icon: "success",
          showConfirmButton: false,
        }).then(() => {
          window.location.href = "/" //Location after deleted is success;
        })
      } else{
        Swal.fire({
          title: 'Error',
          icon: "error",
          text: result.data,
          showConfirmButton: true,
        });
      }
    }
  });
});

$(".edit-profile").show();
$(".account-settings").hide();
