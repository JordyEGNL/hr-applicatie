package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

var isNumber = regexp.MustCompile(`^[0-9]+$`).MatchString
var isPhoneNumber = regexp.MustCompile(`^[0-9-]+$`).MatchString
var isAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9 ]+$`).MatchString

func handleError(c *gin.Context, statusCode int, message string, err error) {
	if err != nil {
		log.Printf("ERROR: %s: %v", message, err)
	} else {
		log.Printf("ERROR: %s", message)
	}
	c.JSON(statusCode, gin.H{"message": message})
}

func validateRequiredFields(employee addEmployeeJson) error {
	if employee.FirstName == "" || employee.LastName == "" || employee.Email == "" || employee.Phone == "" || employee.HireDate == "" || employee.DepartmentID == 0 {
		return errors.New("please provide all required fields")
	}
	return nil
}

func checkDepartmentExists(departmentID int) error {
	departments := getDepartmentFromDB(fmt.Sprint(departmentID))
	if len(departments) == 0 {
		return errors.New("department does not exist")
	}
	return nil
}

func checkEmail(email string) error {
	err := checkUniqueEmail(email)
	if err != nil {
		return err
	}
	err = checkEmailSpaces(email)
	if err != nil {
		return err
	}

	return nil
}

func checkUniqueEmail(email string) error {
	employees := getEmployeeFromDB("")
	for _, emp := range employees {
		if emp.Email == email {
			return errors.New("email already exists")
		}
	}
	return nil
}

func checkEmailSpaces(email string) error {
	if strings.Contains(email, " ") {
		return errors.New("email cannot contain spaces")
	}
	return nil
}

func checkEmailSuffix(email string) error {
	if !strings.HasSuffix(email, "@holiday-parks.eu") {
		return errors.New("email must end with @holiday-parks.eu")
	}
	return nil
}

func checkPhoneNumber(phone string) error {
	if !isPhoneNumber(phone) {
		return errors.New("phone number can only contain numbers and dashes and should not have spaces")
	}
	if len(phone) > 16 {
		return errors.New("phone number is too long")
	}
	return nil
}

func checkName(name string) error {
	// Check if name contains only letters
	if !isAlphaNumeric(name) {
		return errors.New("name can only contain letters and have no spaces")
	}
	// Check if starts with a capital letter
	if name != "" && name[0] != strings.ToUpper(name)[0] {
		return errors.New("name must start with a capital letter")
	}
	// Check max length
	if len(name) > 32 {
		return errors.New("name is too long")
	}
	return nil
}

func checkLastName(lastName string) error {
	// Check if last name contains only letters
	if !isAlphaNumeric(lastName) {
		return errors.New("last name can only contain letters and have no spaces")
	}
	// Check max length
	if len(lastName) > 32 {
		return errors.New("last name is too long")
	}
	return nil
}

func checkAddress(address string) error {
	// Check if address contains only alphanumeric characters
	if !isAlphaNumeric(address) {
		return errors.New("address can only contain alphanumeric characters and have no spaces")
	}
	// Check max length
	if len(address) > 32 {
		return errors.New("address is too long")
	}
	return nil
}
