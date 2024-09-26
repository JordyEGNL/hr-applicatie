async function login() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;
    const errorMessage = document.getElementById('error-message');
    
    // Clear any previous error message
    errorMessage.textContent = '';

    try {
        const response = await fetch('/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: new URLSearchParams({
                'username': username,
                'password': password
            })
        });

        const result = await response.json();

        if (response.status === 200) {
            console.log(result.message);
            // Redirect to frontend
            window.location.href = '/frontend';
        } else {
            errorMessage.textContent = result.error;
        }
    } catch (error) {
        errorMessage.textContent = 'An error occurred during the login process';
        console.error('Error:', error);
    }
}

// Automatically call login on page load
document.addEventListener('DOMContentLoaded', function() {
    login();
});

// login when page loads
window.addEventListener("load", (event) => {
    login();
  });



// Enter key to submit
document.addEventListener('keydown', function(event) {
    if (event.key === 'Enter') {
        login();
    }
});

function viewPassword() {
    var x = document.getElementById("password");
    if (x.type === "password") {
      x.type = "text";
    } else {
      x.type = "password";
    }
  } 