const form = {
    email: $("#email-form"),
    password: $("#password-form"),
    submit: $(".login-signup__btn"),
    pass_alert: $(".pass-alert"),
    email_alert: $(".email-alert"),
    form_message: $(".form-message")
};

let checkForm = () => {
    let isFailed = false;

    // Check if email is empty
    if(form.email.val() === "") {
        form.email_alert.html("Please provide email");
        isFailed = true;
    } else {
        form.email_alert.html("");
    }
    
    // Check if password is empty
    if(form.password.val() === "") {
        form.pass_alert.html("Please provide password");
        isFailed = true;
    } else {
        form.pass_alert.html("");
    }

    return isFailed? false:true;
}

form.submit.on("click",async () => {
    e.preventDefault();
    if(checkForm()) {
        form_message.html("");
        const login = "/login";

        fetch(login, {
            method: "POST",
            redirect: "follow",
            headers: {
                Accept: "application/json, text/plain, */*",
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                username: form.email.val(),
                password: form.password.val(),
            }),
        })
            .then((response) => response.json())
            .then((result) => {
                if (result.error) {
                    form.form_message.html("result.data");
                } else {
                    window.location.href = result.data;
                }
        });
    }
});

