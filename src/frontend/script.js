document.addEventListener("DOMContentLoaded", function() {
    // Fetch current version
    fetch('/healthcheck')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            console.log('Version:', data.version);
            document.getElementById("currentVersion").innerText = data.version;
        })
        .catch(error => console.error('Error fetching current version:', error));

    // Fetch current logged in user
    fetch('/api/v1/me')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            console.log('Logged in user:', data.username);
            document.getElementById("loggedInUser").innerText = data.first_name + " " + data.last_name;
        })
        .catch(error => console.error('Error fetching logged in user:', error));

    // Fetch departments and populate the form and table
    fetch('/api/v1/department')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            console.log('Departments fetched successfully:', data);
            const departmentSelect = document.getElementById('department');
            const tableBody = document.querySelector("#departmentTable tbody");
            data.forEach(department => {
                // Populate select options
                const option = document.createElement('option');
                option.value = department.id;
                option.textContent = department.name + " (" + department.id + ")";
                departmentSelect.appendChild(option);

                // Populate department table
                const row = document.createElement("tr");
                row.innerHTML = `
                    <td>${department.id}</td>
                    <td>${department.name}</td>
                    <td>${department.location}</td>
                    <td>${department.manager}</td>
                `;
                tableBody.appendChild(row);
            });
        })
        .catch(error => {
            console.error('Error fetching departments:', error);
            const errorMessage = document.getElementById('error-message');
            errorMessage.textContent = 'Failed to load departments. Please try again later.';
        });

    // Autofill email field
    const firstNameInput = document.getElementById('firstName');
    const lastNameInput = document.getElementById('lastName');
    const emailInput = document.getElementById('email');

    function updateEmail() {
        const firstName = firstNameInput.value.trim();
        const lastName = lastNameInput.value.trim();

        // Remove spaces and convert to lowercase
        const formattedLastName = lastName.replace(/\s/g, '');

        if (firstName && lastName) {
            const email = `${firstName.charAt(0)}.${formattedLastName}@holiday-parks.eu`.toLowerCase();
            emailInput.value = email;
        }
    }

    firstNameInput.addEventListener('input', updateEmail);
    lastNameInput.addEventListener('input', updateEmail);

    // Update phone number with - after 06 only once
    const phoneNumberInput = document.getElementById('phone');
    phoneNumberInput.addEventListener('input', function() {
        const phoneNumber = phoneNumberInput.value;
        if (phoneNumber.startsWith('06') && !phoneNumber.includes('-')) {
            phoneNumberInput.value = phoneNumber.replace('06', '06-');
        }
    });


    // Fetch and display employees
    function fetchLocations() {
        console.log('Fetching locations...');
        fetch('/api/v1/location', {
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            console.log('Locations fetched successfully:', data);
            const tableBody = document.querySelector("#locationTable tbody");
            tableBody.innerHTML = ''; // Clear existing data
            data.forEach(location => {
                const row = document.createElement("tr");
                row.innerHTML = `
                    <td>${location.id}</td>
                    <td>${location.name}</td>
                    <td>${location.address}</td>
                    <td>${location.city}</td>
                    <td>${location.postal_code}</td>
                    <td>${location.country}</td>
                `;
                tableBody.appendChild(row);
            });
        })
        .catch(error => console.error('Error fetching data:', error));
        }
    fetchLocations();
    fetchDepartments();
});
    
// Logout
function logout() {
    window.location.href = '/logoutpage';
}

function toggleForm() {
    const form = document.getElementById('addEmployeeForm');
    if (form.style.display === 'none' || form.style.display === '') {
        form.style.display = 'grid';
    } else {
        form.style.display = 'none';
    }
}