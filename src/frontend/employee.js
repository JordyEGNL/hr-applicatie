document.addEventListener('DOMContentLoaded', function() {
    const employeeTableBody = document.querySelector('#employeeTable tbody');
    const editEmployeeForm = document.getElementById('editEmployeeForm');
    const editFieldSelect = document.getElementById('editField');
    const editValueInput = document.getElementById('editValue');
    const errorMessage = document.querySelector('.error-message');

    // Focus on search bar when page loads
    document.getElementById('searchBar').focus();

    function resetErrorMessage() {
        errorMessage.textContent = '';
        errorMessage.style.color = 'red';
        errorMessage.style.fontWeight = 'bold';
    }

    /**
     * Fetches employees from the API and populates the table.
     */
    let employeeData = [];
    let filteredEmployeeData = [];
    let currentPage = 0;
    const rowsPerPage = 10;

    function fetchEmployees() {
        // Fetch employees from the API and populate the table
        fetch('/api/v1/employee')
            .then(response => response.json())
            .then(data => {
                employeeData = data;
                filteredEmployeeData = data;
                displayEmployees();
            })
            .catch(error => console.error('Error fetching employees:', error));
    }

    function displayEmployees() {
        employeeTableBody.innerHTML = '';
        const start = currentPage * rowsPerPage;
        const end = start + rowsPerPage;
        const employeesToShow = filteredEmployeeData.slice(start, end);    
        let last_login_modified;
    
        employeesToShow.forEach(employee => {
            last_login_modified = "N/A";
            if (employee.can_login == 1) {
                employee.can_login = '✅';
            } else if (employee.can_login == 0) {
                employee.can_login = '❌';
            }
            if (employee.last_login !== "N/A") {
                // If already UTC time, just convert to local time
                const lastLoginDate = new Date(employee.last_login + " UTC");
                last_login_modified = lastLoginDate.toLocaleString();
                if (lastLoginDate.toDateString() === new Date().toDateString()) {
                    last_login_modified = lastLoginDate.toLocaleTimeString();
                }
                if (new Date() - lastLoginDate < 60000) {
                    last_login_modified = 'Just now';
                } else if (new Date() - lastLoginDate >= 60000 && new Date() - lastLoginDate < 3600000) {
                    last_login_modified = Math.floor((new Date() - lastLoginDate) / 60000) + ' minutes ago';
                }
            } else if (employee.last_login === "N/A") {
                last_login_modified = "➖";
            }
    
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${employee.id}</td>
                <td>${employee.first_name}</td>
                <td>${employee.last_name}</td>
                <td>${employee.email}</td>
                <td>${employee.phone}</td>
                <td>${employee.hire_date}</td>
                <td>${employee.department}</td>
                <td>${employee.manager_email}</td>
                <td>${employee.location.name}</td>
                <td>${employee.can_login}</td>
                <td>${last_login_modified}</td>
                <td>
                    <button class="edit-button databutton" data-id="${employee.id}">✏️</button>
                    <button class="delete-button databutton" data-id="${employee.id}">🗑️</button>
                </td>
            `;
            employeeTableBody.appendChild(row);
        });
    
        // Add event listeners for edit buttons
        document.querySelectorAll('.edit-button').forEach(button => {
            button.addEventListener('click', handleEditButtonClick);
        });
    
        // Add event listeners for delete buttons
        document.querySelectorAll('.delete-button').forEach(button => {
            button.addEventListener('click', handleDeleteButtonClick);
        });
    
        // Show or hide the "Show More" and "Show Previous" buttons
        const showMoreButton = document.getElementById('showMoreButton');
        const showPreviousButton = document.getElementById('showPreviousButton');
        if (end < filteredEmployeeData.length) {
            showMoreButton.classList.remove('hidden');
            showMoreButton.disabled = false;
        } else {
            showMoreButton.classList.add('hidden');
            showMoreButton.disabled = true;
        }
        if (currentPage > 0) {
            showPreviousButton.classList.remove('hidden');
            showPreviousButton.disabled = false;
        } else {
            showPreviousButton.classList.add('hidden');
            showPreviousButton.disabled = true;
        }

        // Generate page numbers
        const pageNumbersContainer = document.getElementById('pageNumbers');
        pageNumbersContainer.innerHTML = '';
        const totalPages = Math.ceil(filteredEmployeeData.length / rowsPerPage);
        for (let i = 0; i < totalPages; i++) {
            const pageNumberButton = document.createElement('button');
            pageNumberButton.textContent = i + 1;
            pageNumberButton.addEventListener('click', () => goToPage(i));
            if (i === currentPage) {
                pageNumberButton.classList.add('active');
            }
            pageNumbersContainer.appendChild(pageNumberButton);
        }
    }
    
    function showNextEmployees() {
        if (currentPage * rowsPerPage < filteredEmployeeData.length - rowsPerPage) {
            currentPage++;
            displayEmployees();
        }
    }

    function showPreviousEmployees() {
        if (currentPage > 0) {
            currentPage--;
            displayEmployees();
        }
    }

    function goToPage(page) {
        currentPage = page;
        displayEmployees();
    }
    
    function searchEmployees() {
        const query = document.getElementById('searchBar').value.toLowerCase();
        filteredEmployeeData = employeeData.filter(employee => {
            return (
                employee.first_name.toLowerCase().includes(query) ||
                employee.last_name.toLowerCase().includes(query) ||
                employee.email.toLowerCase().includes(query) ||
                employee.phone.toLowerCase().includes(query) ||
                employee.department.toLowerCase().includes(query) ||
                employee.manager_email.toLowerCase().includes(query) ||
                employee.location.name.toLowerCase().includes(query)
            );
        });
        currentPage = 0; // Reset to the first page
        displayEmployees();
    }
    
    document.getElementById('showMoreButton').addEventListener('click', showNextEmployees);
    document.getElementById('showPreviousButton').addEventListener('click', showPreviousEmployees);
    document.getElementById('searchBar').addEventListener('input', searchEmployees);
    
    // Initialize the fetch
    fetchEmployees();

    // When a user clicks the edit button, show the edit form
    function handleEditButtonClick(event) {
        const employeeId = event.target.getAttribute('data-id');
        editEmployeeForm.employeeId.value = employeeId;
        editEmployeeForm.style.display = 'grid';

        // move the focus to the first input field
        editValueInput.focus();
    
        // Find the employee data from the employeeData array
        const employee = employeeData.find(emp => emp.id == employeeId);

        if (employee) {
            const selectedField = editFieldSelect.value;
            const currentValue = employee[selectedField];
            editValueInput.value = ''; // Clear the value
            editValueInput.placeholder = currentValue; // Set placeholder as current value

            // Update the value when the field changes
            editFieldSelect.addEventListener('change', () => {
                const selectedField = editFieldSelect.value;
                const currentValue = employee[selectedField];
                editValueInput.value = ''; // Clear the value
                if (currentValue == "✅" || currentValue == "❌") {
                    editValueInput.checked = currentValue == "✅" ? true : false;
                }
                editValueInput.placeholder = currentValue; // Set placeholder as current value
            });
        } else {
            console.error('Employee not found in employeeData array');
        }
    }

    function handleDeleteButtonClick(event) {
        const employeeId = event.target.getAttribute('data-id');
        console.log(`Delete employee ${employeeId}`);
        fetch(`/api/v1/employee/${employeeId}`, {
            method: 'DELETE'
        })
        .then(response => {
            resetErrorMessage();
            if (!response.ok) {
                return response.json().then(errorData => {
                    errorMessage.textContent = "Kan gebruiker niet verwijderen: " + errorData.message;
                    throw new Error(errorData.message);
                });
            }
            errorMessage.style.color = 'green';
            errorMessage.textContent = 'Gebruiker is verwijderd.';
            fetchEmployees();
            return response.json();
        })
    }

    editEmployeeForm.addEventListener('submit', function(event) {
        event.preventDefault();

        const formData = new FormData(editEmployeeForm);
        const employeeId = formData.get('employeeId');
        const editField = formData.get('editField');
        let editValue = formData.get('editValue');

        // Handle specific data types
        if (editField === 'department_id') {
            editValue = parseInt(editValue);
        } else if (editField === 'can_login') {
            editValue = editValue === "true" ? "true" : "false";
        }

        const updateData = {
            [editField]: editValue
        };

        console.log(`Editing employee ${employeeId}:`, updateData);

        fetch(`/api/v1/employee/${employeeId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(updateData)
        })
        .then(response => {
            resetErrorMessage();
            if (!response.ok) {
                return response.json().then(errorData => {
                    errorMessage.textContent = "Kan gebruiker niet bewerken: " + errorData.message;
                    throw new Error(errorData.message);
                });
            }
            return response.json();
        })
        .then(data => {
            if (data.message === 'User updated') {
                console.log('Employee updated successfully');
                fetchEmployees();
                editEmployeeForm.reset();
                editEmployeeForm.style.display = 'none';
                errorMessage.style.color = 'green';
                errorMessage.textContent = 'Gebruiker is bijgewerkt.';
            } else {
                errorMessage.textContent = data.error;
                alert(data.message);
            }
        })
        .catch(error => console.error('Error editing employee:', error));
    });

    // Handle form submission for adding a new employee
    const addEmployeeForm = document.getElementById('addEmployeeForm');
    addEmployeeForm.addEventListener('submit', function(event) {
        event.preventDefault();

        const formData = new FormData(addEmployeeForm);
        const employeeData = {
            first_name: formData.get('firstName'),
            last_name: formData.get('lastName'),
            email: formData.get('email'),
            phone: formData.get('phone'),
            hire_date: formData.get('hireDate'),
            department_id: parseInt(formData.get('department')),
            can_login: formData.get('canLogin') === "on" ? "true" : "false"
        };

        console.log('Adding new employee:', employeeData);

        fetch('/api/v1/employee', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(employeeData)
        })
        .then(response => {
            resetErrorMessage();
            // if response is not 200 OK, show the error response from the server
            if (!response.ok) {
                return response.json().then(errorData => {
                    errorMessage.textContent = "Kan gebruiker niet toevoegen: " + errorData.message;
                    throw new Error(errorData.message);
                });
            }
            return response.json();
        })
        .then(data => {
            if (data.message == 'User added'){
                console.log('Employee added successfully');
                fetchEmployees();
                addEmployeeForm.reset();
                errorMessage.style.color = 'green';
                errorMessage.textContent = 'Gebruiker ' + employeeData.first_name + ' ' + employeeData.last_name + ' is toegevoegd.';
                addEmployeeForm.style.display = 'none';
            } else {
                errorMessage.textContent = result.error;
                alert(data.message);
            }
        })
        .catch(error => console.error('Error adding employee:', error));
    });
});
