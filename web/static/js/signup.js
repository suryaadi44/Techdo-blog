const form = {
    email: $("#email-form"),
    username: $("#username-form"),
    firstName: $("#first-name-form"),
    lastName: $("#last-name-form"),
    password: $("#password-form"),
    confirm_pass: $("#retype-password-form"),
}

const alert = {
    emailAlert: $(".email-alert"),
    usernameAlert: $(".username-alert"),
    firstNameAlert: $(".first-name-alert"),
    lastNameAlert: $(".last-name-alert"),
    passwordAlert: $(".password-alert"),
    retypeAlert: $(".retype-password-alert"),
    passwordMatchAlert: $(".password-match-alert"),
    signupMessage: $(".signup__message-alert")
}

const signupBtn = $(".login-signup__btn");

const checkForm = () => {
    let isFailed = false;

    if(form.email.val() === "") {
        alert.emailAlert.html("*Please provide email");
        isFailed = true 
    } else {
        alert.emailAlert.html("");
    }

    if(form.username.val() === "") {
        alert.usernameAlert.html("*Please provide username");
        isFailed = true 
    } else {
        alert.usernameAlert.html("");
    }

    if(form.firstName.val() === "") {
        alert.firstNameAlert.html("*Please provide first name");
        isFailed = true 
    } else {
        alert.firstNameAlert.html("");
    }

    if(form.lastName.val() === "") {
        alert.lastNameAlert.html("*Please provide last name");
        isFailed = true 
    } else {
        alert.lastNameAlert.html("");
    }

    if(form.password.val() === "") {
        alert.passwordAlert.html("*Please provide password");
        isFailed = true 
    } else {
        alert.passwordAlert.html("");
    }

    if(form.confirm_pass.val() === "") {
        alert.retypeAlert.html("*Please provide re-type password");
        isFailed = true 
    } else {
        alert.retypeAlert.html("");
    }

    return isFailed? false:true; 
}

form.confirm_pass.on("focusout keyup", () => {
    if (form.confirm_pass.val() === "") {
        alert.passwordMatchAlert.html("");
        isFailed = true 
    } else {
        alert.retypeAlert.html("");
        if (form.password.val() != form.confirm_pass.val()) {
            alert.passwordMatchAlert.css("color", "red");
            alert.passwordMatchAlert.html("*Password is not match")
            isFailed = true 
        } else {
            alert.passwordMatchAlert.css("color", "green");
            alert.passwordMatchAlert.html("*Password match")
            filled += 1;
        }
    }

});

signupBtn.on("click", () => {
    alert.signupMessage.html("");
    if (checkForm()) {
        const signup = "/signup";
        fetch(signup, {
            method: "POST",
            redirect: "follow",
            headers: {
                Accept: "application/json, text/plain, */*",
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                email: form.username.val(),
                username: form.username.val(),
                firstName: form.firstName.val(),
                lastName: form.lastName.val(),
                password: form.password.val(),
            }),
        })
            .then((response) => response.json())
            .then((result) => {
                if (result.error) {
                    alert.signupMessage.html(result.data);
                } else {
                    window.location.href = result.data;
                }
            });
    } 

});