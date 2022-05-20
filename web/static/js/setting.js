const navBtn = {
  editProfileNav: $(".edit-profile-nav"),
  accountSettingNav: $(".account-setting-nav")
};

async function getUserMiniDetail() {
  try {
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

async function deleteAccount() {
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

navBtn.editProfileNav.on("click", () => {
  $(".edit-profile").show();
  $(".account-settings").hide();

  const form = {
    profile: document.querySelector("#profile-pic"),
    firstName: $("#first-name-form"),
    lastName: $("#last-name-form"),
    phone: $("#phone-form"),
    aboutMe: $("#about-form"),
    submit: $("#submit-btn"),
    submitPicture: $("#submit-picture"),
  };

  form.submit.on("click", async (e) => {
    e.preventDefault();
    try {
      let res = await fetch("/user/detail", {
        method: 'POST',
        credentials: 'include',
        body: JSON.stringify({
          first_name: form.firstName.val(),
          last_name: form.lastName.val(),
          phone: form.phone.val(),
          about_me: form.aboutMe.val(),
        }),
      });
      let response = await res.json();

      if (!response.error) {
        Swal.fire({
          title: 'Account Updated',
          icon: "success",
        }).then(function () {
          location.reload()
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
  

  form.submitPicture.on("click", async (e) => {
    const formData = new FormData();

    formData.append('profile-pic', form.profile.files[0]);
    try {
      let res = await fetch("/user/detail/picture", {
        method: 'POST',
        credentials: 'include',
        body: formData,
      });
      let response = await res.json();

      if (!response.error) {
        Swal.fire({
          title: 'Picture Updated',
          icon: "success",
        }).then(function () {
          location.reload()
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
      } else {
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

navBtn.editProfileNav.click()