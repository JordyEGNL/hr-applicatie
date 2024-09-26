function loginScreen() {
    window.location.href = '/login';
}

// on page load run logout()
window.addEventListener("load", (event) => {
    logout();
  });

// removes cookies
async function logout() {
    fetch('/logout', {
        method: 'GET'}
    )
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .catch(error => console.error('Error logging out:', error));
}