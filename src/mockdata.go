package main

import (
	"github.com/gin-gonic/gin"
)

// mock data for testing
func mockData(c *gin.Context) {
	// Add locations
	addLocationToDB(addLocationJson{
		Name:       "Location 1",
		Address:    "Address 1",
		PostalCode: "1234",
		City:       "City 1",
		Country:    "Country 1",
	})
	addLocationToDB(addLocationJson{
		Name:       "Location 2",
		Address:    "Address 2",
		PostalCode: "5678",
		City:       "City 2",
		Country:    "Country 2",
	})

	// Add employees
	addEmployeeToDB(addEmployeeJson{
		FirstName:    "Employee 1",
		LastName:     "Lastname 1",
		Email:        "test@holiday-parks.eu",
		Phone:        "1234567890",
		HireDate:     "2021-01-01",
		DepartmentID: 1,
	})
	addEmployeeToDB(addEmployeeJson{
		FirstName:    "Employee 2",
		LastName:     "Lastname 2",
		Email:        "test2@holiday-parks.eu",
		Phone:        "0987654321",
		HireDate:     "2021-01-01",
		DepartmentID: 2,
	})
	addEmployeeToDB(addEmployeeJson{
		FirstName:    "Admin 1",
		LastName:     "Lastname 1",
		Email:        "admin@holiday-parks.eu",
		Phone:        "1234567890",
		HireDate:     "2021-01-01",
		DepartmentID: 1,
		CanLogin:     "true",
		Password:     "Admin01!",
	})

	// Add departments
	addDepartmentToDB(addDepartmentJson{
		Name:       "Department 1",
		LocationID: 1,
		ManagerID:  1,
	})
	addDepartmentToDB(addDepartmentJson{
		Name:       "Department 2",
		LocationID: 2,
		ManagerID:  2,
	})
}
