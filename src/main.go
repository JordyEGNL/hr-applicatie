package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "holiday-parks.eu/hr-api/docs"
)

const currentVersion string = "2024.06.14.a"

//	@title			HR API
//	@description	This is the HR API for the Holiday Parks company.

//	@contact.name	Jordy Hoebergen
//	@contact.email	j.hoebergen@student.fontys.nl

//	@license.name	GPL-3.0
//	@license.url	https://www.gnu.org/licenses/gpl-3.0.en.html

//	@host		127.0.0.1:5000
//	@BasePath	/
//  @schemes	http

var db *sql.DB

func init() {
	// Log to a file
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	writer := io.MultiWriter(os.Stdout, file)
	log.SetOutput(writer)

	// Get env from docker and set it to the config
	getEnv()

	// Print the version
	log.Printf("INFO: Starting API version %s", currentVersion)

	// Print debug information
	if api.debug {
		log.Printf("")
		log.Printf("------------------------ DEBUG ------------------------")
		log.Printf("Database credentials:")
		log.Printf("  Host: %s\n", dbcred.host)
		log.Printf("  Port: %s\n", dbcred.port)
		log.Printf("  User: %s\n", dbcred.user)
		log.Printf("  Password: %s\n", dbcred.password)

		log.Printf("API configuration:")
		log.Printf("  Default user: %s\n", api.defaultuser)
		log.Printf("  Default password: %s\n", api.defaultpassword)
		log.Printf("  Debug: %v\n", api.debug)
		log.Printf("------------------------ DEBUG ------------------------")
		log.Printf("")
	}

	if !api.debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to the database
	// Keep trying to connect until it is successful
	err = initializeDB()
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	createTables()

}

func main() {
	log.Printf("INFO: Starting the API on port %v", api.port)

	// Create a new instance of the gin router
	r := gin.Default()

	// Setup the cookie store for session management
	r.Use(sessions.Sessions("mysession", cookie.NewStore(secret)))

	// Login and logout routes
	r.POST("/login", login)
	r.GET("/logout", logout)

	// Private group, require authentication to access
	private := r.Group("/")
	private.Use(AuthRequired)
	{
		private.GET("/api/v1/me", me)

		// Employee endpoints
		private.GET("/api/v1/employee", getEmployee)           // All employees
		private.GET("/api/v1/employee/:id", getEmployee)       // Specific employee
		private.POST("/api/v1/employee", addEmployee)          // Add a employee
		private.PUT("/api/v1/employee/:id", updateEmployee)    // Update a employee
		private.DELETE("/api/v1/employee/:id", deleteEmployee) // Delete a employee

		// Department endpoints
		private.GET("/api/v1/department", getDepartment)           // All departments
		private.GET("/api/v1/department/:id", getDepartment)       // Specific department
		private.POST("/api/v1/department", addDepartment)          // Add a department
		private.PUT("/api/v1/department/:id", updateDepartment)    // Update a department
		private.DELETE("/api/v1/department/:id", deleteDepartment) // Delete a department

		// Endpoint for location
		private.GET("/api/v1/location", getLocation)           // All locations
		private.GET("/api/v1/location/:id", getLocation)       // Specific location
		private.POST("/api/v1/location", addLocation)          // Add a location
		private.PUT("/api/v1/location/:id", updateLocation)    // Update a location
		private.DELETE("/api/v1/location/:id", deleteLocation) // Delete a location

		private.GET("/api/v1/mockData", mockData) // Inject mock data

		private.Static("/frontend", "./frontend")
		private.Static("/logoutpage", "./frontend/logout")
	}

	r.GET("/healthcheck", healthCheck) // Health check

	// Login page
	r.Static("/login", "./frontend/login")

	// Redirect / to /login
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/login")
	})

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run("0.0.0.0:" + api.port)
}

// Function to connect to the database
func initializeDB() error {
	cfg := mysql.Config{
		User:                 dbcred.user,
		Passwd:               dbcred.password,
		Net:                  "tcp",
		Addr:                 dbcred.host + ":" + dbcred.port,
		AllowNativePasswords: true,
	}

	// FormatDSN will return a DSN string that
	// can be used to open a database connection
	obfuscatedConnectionString := obfuscatePassword(cfg.FormatDSN())
	if api.debug {
		log.Printf("DEBUG: Connection string: " + obfuscatedConnectionString)
	}

	var err error
	// Retry logic
	for {
		// Connect to the MySQL server without specifying a database
		db, err = sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			log.Printf("ERROR: Cannot connect to the MySQL server: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Check if the database connection is established
		err = db.Ping()
		if err != nil {
			log.Printf("ERROR: Cannot ping the database: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}

	// Check if the "hr" database exists
	query := "CREATE DATABASE IF NOT EXISTS " + dbcred.dbname
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	db.Close()

	cfg.DBName = dbcred.dbname

	// Connect to the MySQL server and "hr" database directly
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return err
	}

	// Set the maximum number of open connections to the database
	db.SetMaxOpenConns(25)
	// Set the maximum number of idle connections in the pool
	db.SetMaxIdleConns(25)
	// Set the maximum lifetime of a connection
	db.SetConnMaxLifetime(time.Minute * 5)

	// Check if the database connection is established
	err = db.Ping()
	if err != nil {
		return err
	}

	if api.debug {
		log.Printf("DEBUG: Connected to the database")
	}
	return nil
}

func obfuscatePassword(connectionString string) string {
	// Split the connection string by the colons
	// The password is at the first index.
	// Example: user:password@tcp(host:port)/database
	parts := strings.Split(connectionString, ":")
	if len(parts) < 2 {
		return connectionString
	}

	// Get the index of the password
	// The password ends at the index of the '@' symbol
	passwordEnd := strings.Index(parts[1], "@")
	if passwordEnd == -1 {
		return connectionString
	}

	// Get the password from the connection string
	password := parts[1][:passwordEnd]
	obfuscatedPassword := strings.Repeat("*", len(password))
	parts[1] = strings.Replace(parts[1], password, obfuscatedPassword, 1)

	return strings.Join(parts, ":")
}

func createTables() {
	// Create the employee table
	// Create the table if it does not exist
	query := `CREATE TABLE IF NOT EXISTS employee (
		id INT AUTO_INCREMENT PRIMARY KEY,
		first_name VARCHAR(100) DEFAULT '',
		last_name VARCHAR(100) DEFAULT '',
		email VARCHAR(100) DEFAULT '',
		phone VARCHAR(20) DEFAULT '',
		hire_date DATE DEFAULT '1970-01-01',
		department_id INT DEFAULT '0',
		can_login BOOLEAN DEFAULT FALSE,
		password VARCHAR(100) DEFAULT '',
		creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_login TIMESTAMP NULL DEFAULT NULL,
		edit_date TIMESTAMP NULL DEFAULT NULL
	)`
	_, err := db.Query(query)
	if err != nil {
		log.Fatalf("ERROR: Cannot create table employee: %v", err)
	}

	// Create the department table
	// Create the table if it does not exist
	query = `CREATE TABLE IF NOT EXISTS department (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100),
		location_id INT,
		manager_id INT
	)`
	_, err = db.Query(query)
	if err != nil {
		log.Fatalf("ERROR: Cannot create table department: %v", err)
	}

	// Create the location table
	// Create the table if it does not exist
	query = `CREATE TABLE IF NOT EXISTS location (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100),
		address VARCHAR(100),
		postal_code VARCHAR(10),
		city VARCHAR(100),
		country VARCHAR(100)
	)`
	_, err = db.Query(query)
	if err != nil {
		log.Fatalf("ERROR: Cannot create table location: %v", err)
	}

	// Add foreign keys after tables are created
	query = `ALTER TABLE employee ADD CONSTRAINT fk_department FOREIGN KEY (department_id) REFERENCES department(id)`
	_, err = db.Query(query)
	if err != nil {
		log.Fatalf("ERROR: Cannot add foreign key to employee: %v", err)
	}

	query = `ALTER TABLE department ADD CONSTRAINT fk_location FOREIGN KEY (location_id) REFERENCES location(id)`
	_, err = db.Query(query)
	if err != nil {
		log.Fatalf("ERROR: Cannot add foreign key to department: %v", err)
	}

	query = `ALTER TABLE department ADD CONSTRAINT fk_manager FOREIGN KEY (manager_id) REFERENCES employee(id)`
	_, err = db.Query(query)
	if err != nil {
		log.Fatalf("ERROR: Cannot add foreign key to department: %v", err)
	}

	// If no users are in the database, add a default user
	query = `SELECT COUNT(*) FROM employee`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("ERROR: Cannot count employees: %v", err)
	}

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			log.Fatalf("ERROR: Cannot scan row: %v", err)
		}
	}

	if count == 0 {
		// Add a default user
		query = `INSERT INTO employee (first_name, last_name, email, phone, hire_date, department_id, can_login, password) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
		_, err = db.Query(query, "Luka", "Van Zutven", api.defaultuser, "0612345678", "2021-01-01", 1, true, hashPassword(api.defaultpassword))
		if err != nil {
			log.Fatalf("ERROR: Cannot insert default user: %v", err)
		}
		log.Printf("INFO: Default user added")
	}
}

const userkey = "user"

var secret = []byte("secret")

// StatusBadRequest structure for 400 response
type StatusBadRequest struct {
	Error string `json:"error"`
}

// StatusUnauthorized structure for 401 response
type StatusUnauthorized struct {
	Message string `json:"message"`
}

// StatusOK structure for 200 response
type StatusOK struct {
	Message string `json:"message"`
}

// StatusInternalServerError structure for 500 response
type StatusInternalServerError struct {
	Error string `json:"error"`
}

// AuthRequired is a simple middleware to check the session.
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}
	// Continue down the chain to handler etc

	// Log the login attempt
	log.Printf("INFO: User %s is authenticated from ip %v accessing %v", user, c.ClientIP(), c.Request.URL.Path)
	c.Next()
}

func hashPassword(password string) string {
	if password == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(password))

	// Return the hash as a string
	return fmt.Sprintf("%x", hash)
}

// Login
//
//	@Summary		Authenticate user
//	@Description	Autenhticate the user with the provided credentials to get access to the API and frontend
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200			{object}	StatusOK
//	@Failure		401			{object}	StatusUnauthorized
//	@Param			username	header		string	false	"Username"
//	@Param			password	header		string	false	"Password"
//	@Router			/login [post]
func login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.GetHeader("username")
	// Hash the password from the form
	password := hashPassword(c.GetHeader("password"))

	// If header is empty use form input
	if username == "" || password == "" {
		username = c.PostForm("username")
		password = hashPassword(c.PostForm("password"))
	}

	// if cookie is still valid return 200
	if session.Get(userkey) != nil {
		c.JSON(http.StatusOK, gin.H{"message": "User is already authenticated"})
		return
	}

	// Convert hash to lowercase
	password = strings.ToLower(password)

	// Convert email from @fictproftaak03.onmicrosoft.com to @holiday-parks.eu to match the database
	username = strings.Replace(username, "@fictproftaak03.onmicrosoft.com", "@holiday-parks.eu", 1)

	// Log the login attempt
	log.Printf("INFO: User %s is trying to login with hash %s from ip %s", username, password, c.ClientIP())

	// Print all the received headers in the console for debugging
	if api.debug {
		log.Printf("----------- Headers -----------")
		for name, values := range c.Request.Header {
			for _, value := range values {
				log.Printf("%s: %s", name, value)
			}
		}
		log.Printf("----------- Headers -----------")
	}

	// Check if user is allowed to login in db
	query := `SELECT can_login, email, password FROM employee WHERE email = ? AND password = ?`
	rows, err := db.Query(query, username, password)
	if err != nil {
		log.Printf("ERROR: Cannot get can_login: %v", err)
	}

	var canLogin bool
	var db_email, db_password string
	for rows.Next() {
		err := rows.Scan(&canLogin, &db_email, &db_password)
		if err != nil {
			log.Printf("ERROR: Cannot scan row: %v", err)
		}
	}

	if db_email != username || db_password != password {
		// Send correct password to error message
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})

		// Log the failed login attempt
		log.Printf("INFO: Failed login attempt for user %s with password %s", username, password)
		return
	}

	if !canLogin {
		// Log the failed login attempt
		log.Printf("INFO: Failed login attempt for user %s with password %s", username, password)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User is not allowed to login"})
		return
	}

	// Save the username in the session
	session.Set(userkey, username)
	session.Options(sessions.Options{
		MaxAge:   3600, // 1 hour in seconds
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	// Log the successful login attempt
	log.Printf("INFO: User %s successfully logged in from %v", username, c.ClientIP())
	logUserLogin(username)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
}

// Logout
//
//	@Summary		Logout user
//	@Description	Logout the user and remove the session
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	StatusOK
//	@Failure		400	{object}	StatusBadRequest
//	@Failure		401	{object}	StatusUnauthorized
//	@Failure		500	{object}	StatusInternalServerError
//	@Router			/logout [get]
func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

func logUserLogin(username string) {
	// Log the login attempt
	query := `UPDATE employee SET last_login = CURRENT_TIMESTAMP WHERE email = ?`
	_, err := db.Query(query, username)
	if err != nil {
		log.Printf("ERROR: Cannot log user login: %v", err)
	}
}

// me is the handler that will return the user information stored in the session.
func me(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get(userkey)

	// Log the login attempt
	query := `SELECT first_name, last_name, last_login FROM employee WHERE email = ?`
	rows, err := db.Query(query, username)
	if err != nil {
		log.Printf("ERROR: Cannot log user login: %v", err)
	}

	var first_name, last_name, last_login string

	for rows.Next() {
		err := rows.Scan(&first_name, &last_name, &last_login)
		if err != nil {
			log.Printf("ERROR: Cannot scan row: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{"username": username,
			"first_name": first_name,
			"last_name":  last_name,
			"last_login": last_login})

	}

}

// Health check
func healthCheck(c *gin.Context) {
	err := db.Ping()
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Database is not reachable",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "OK",
		"version": currentVersion,
	})
}

//
// Employees
//

// struct for the JSON body of the addEmployee endpoint
type addEmployeeJson struct {
	FirstName    string `json:"first_name" example:"John"`
	LastName     string `json:"last_name" example:"Doe"`
	Email        string `json:"email" example:"j.doe@holiday-parks.eu"`
	Phone        string `json:"phone" example:"06-12345678"`
	HireDate     string `json:"hire_date" example:"2021-01-01"`
	DepartmentID int    `json:"department_id" example:"1"`
	CanLogin     string `json:"can_login" example:"true"`
	Password     string `json:"password" example:"Admin01!"`
}

// struct for the JSON body of the getEmployee endpoint
type getEmployeeJson struct {
	ID           int             `json:"id" example:"1"`
	FirstName    string          `json:"first_name" example:"John"`
	LastName     string          `json:"last_name" example:"Doe"`
	Email        string          `json:"email" example:"j.doe@holiday-parks.eu"`
	Phone        string          `json:"phone" example:"06-12345678"`
	HireDate     string          `json:"hire_date" example:"2021-01-01"`
	Department   string          `json:"department" example:"1"`
	ManagerEmail string          `json:"manager_email" example:"m.anager@holiday-parks.eu"`
	Location     getLocationJson `json:"location"`
	CanLogin     bool            `json:"can_login" example:"true"`
	CreationDate string          `json:"creation_date" example:"2021-01-01 9:00:00"`
	EditDate     string          `json:"edit_date" example:"2021-01-01 12:00:00"`
	LastLogin    string          `json:"last_login" example:"2021-01-01 10:00:00"`
}

// GetEmployee
//
//	@Summary		Show employee(s)
//	@Description	Get all information from all employees or a specific employee
//	@Tags			employees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	false	"Employee ID" example:"1"
//	@Success		200	{object}	getEmployeeJson
//	@Failure		401	{object}	StatusUnauthorized
//	@Security		apikey
//	@Router			/api/v1/employee/{id} [get]
func getEmployee(c *gin.Context) {
	// If ID is provided, return the employee with that ID
	id := c.Param("id")
	if id != "" {
		results := getEmployeeFromDB(id)
		c.JSON(http.StatusOK, results)
	} else {
		results := getEmployeeFromDB("")
		c.JSON(http.StatusOK, results)
	}
}

// getEmployeeFromDB gets the employee(s) from the database
func getEmployeeFromDB(givenID string) []getEmployeeJson {
	rows := &sql.Rows{}
	var employees []getEmployeeJson
	var err error

	if givenID != "" {
		// Get the employee with the provided ID
		// Inner join with department and manager
		query := `SELECT 
		e.id, e.first_name, e.last_name, e.email, e.phone, e.hire_date, d.name as department, 
		m.email as manager_email, l.id as location_id, l.name as location_name, 
		l.address, l.postal_code, l.city, l.country as location_country, 
		e.can_login, e.creation_date, e.edit_date, e.last_login 
	  FROM employee e
	  INNER JOIN department d ON e.department_id = d.id
	  INNER JOIN employee m ON d.manager_id = m.id
	  INNER JOIN location l ON d.location_id = l.id
	  WHERE e.id = ?`
		rows, err = db.Query(query, givenID)
		if err != nil {
			log.Printf("ERROR: Cannot get individual employee data: %v", err)
		}
	} else {
		// Get all employees
		// Inner join with department and manager
		query := `SELECT 
		e.id, e.first_name, e.last_name, e.email, e.phone, e.hire_date, d.name as department, 
		m.email as manager_email, l.id as location_id, l.name as location_name, 
		l.address, l.postal_code, l.city, l.country as location_country, 
		e.can_login, e.creation_date, e.edit_date, e.last_login 
	  FROM employee e
	  INNER JOIN department d ON e.department_id = d.id
	  INNER JOIN employee m ON d.manager_id = m.id
	  INNER JOIN location l ON d.location_id = l.id
	  ORDER BY e.id`
		rows, err = db.Query(query)
		if err != nil {
			log.Printf("ERROR: Cannot get employee data: %v", err)
		}
	}

	// If no rows are returned, return none
	if rows == nil {
		return nil
	}

	// Create a list of employees
	var employee getEmployeeJson
	var last_login, edit_date sql.NullString
	var location getLocationJson

	for rows.Next() {
		err := rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.Email, &employee.Phone,
			&employee.HireDate, &employee.Department, &employee.ManagerEmail,
			&location.ID, &location.Name, &location.Address, &location.PostalCode,
			&location.City, &location.Country,
			&employee.CanLogin, &employee.CreationDate, &edit_date, &last_login)
		if err != nil {
			log.Printf("ERROR: Cannot scan row: %v", err)
		}

		// Convert sql.NullString to string
		if last_login.Valid {
			employee.LastLogin = last_login.String
		} else {
			employee.LastLogin = "N/A"
		}

		if edit_date.Valid {
			employee.EditDate = edit_date.String
		} else {
			employee.EditDate = "N/A"
		}

		// if phone is null set it to N/A
		if employee.Phone == "" {
			employee.Phone = "N/A"
		}

		// Create a employee object map values to the json struct
		employee.Location = location

		employees = append(employees, employee)

	}

	// return as json
	return employees
}

// AddEmployee
//
//	@Summary		Add employee
//	@Description	Add a new employee to the database
//	@Tags			employees
//	@Accept			json
//	@Produce		json
//	@Param			first_name		body		string	true	"First name"
//	@Param			last_name		body		string	true	"Last name"
//	@Param			email			body		string	true	"Email"
//	@Param			phone			body		string	true	"Phone"
//	@Param			hire_date		body		string	true	"Hire date" default "1970-01-01"
//	@Param			department_id	body		int		true	"Department ID"
//	@Param			can_login		body		string	true	"Can login" default "false"
//	@Success		200				{object} 	StatusOK
//	@Failure		400				{object}	StatusBadRequest
//	@Failure		401				{object}	StatusUnauthorized
//	@Security		apikey
//	@Router			/api/v1/employee [post]
func addEmployee(c *gin.Context) {
	var employee addEmployeeJson
	err := c.BindJSON(&employee)
	if err != nil {
		handleError(c, 400, "Please provide a valid JSON", err)
		return
	}

	// Validate required fields
	if err := validateRequiredFields(employee); err != nil {
		handleError(c, 400, err.Error(), nil)
		return
	}

	// Check if the department exists
	if err := checkDepartmentExists(employee.DepartmentID); err != nil {
		handleError(c, 400, err.Error(), err)
		return
	}

	// Check email
	if err := checkEmail(employee.Email); err != nil {
		handleError(c, 400, err.Error(), err)
		return
	}

	// Check if the phone number contains only numbers
	if err := checkPhoneNumber(employee.Phone); err != nil {
		handleError(c, 400, err.Error(), err)
		return
	}

	// Check if name contains only letters
	if err := checkName(employee.FirstName); err != nil {
		handleError(c, 400, err.Error(), err)
		return
	}

	if err := checkLastName(employee.LastName); err != nil {
		handleError(c, 400, err.Error(), err)
		return
	}

	newID := addEmployeeToDB(employee)
	c.JSON(200, gin.H{
		"message": "User added",
		"id":      newID,
	})
}

// does the actual database query
func addEmployeeToDB(employee addEmployeeJson) int {
	var err error

	// Hash the password
	employee.Password = hashPassword(api.defaultpassword)

	if employee.CanLogin == "true" {
		employee.CanLogin = "1"
	} else {
		employee.CanLogin = "0"
	}

	// Insert the employee into the database
	query := `INSERT INTO employee (first_name, last_name, email, phone, hire_date, department_id, can_login, password) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = db.Query(query, employee.FirstName, employee.LastName, employee.Email, employee.Phone, employee.HireDate, employee.DepartmentID, employee.CanLogin, employee.Password)
	if err != nil {
		log.Printf("ERROR: Cannot insert employee: %v", err)
	}

	// Get the ID of the inserted employee
	query = `SELECT id FROM employee WHERE first_name = ? AND last_name = ? AND email = ? AND phone = ? AND hire_date = ? AND department_id = ?`
	rows, err := db.Query(query, employee.FirstName, employee.LastName, employee.Email, employee.Phone, employee.HireDate, employee.DepartmentID)
	if err != nil {
		log.Printf("ERROR: Cannot get ID: %v", err)
	}

	// Get the ID from the row
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Printf("ERROR: Cannot scan row: %v", err)
		}
	}

	// Send the ID back with the return statement
	return id
}

// UpdateEmployee
//
//	@Summary		Update employee
//	@Description	Update an employee in the database
//	@Tags			employees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"Employee ID"
//	@Param			first_name	body		string	false	"First name"
//	@Param			last_name	body		string	false	"Last name"
//	@Param			email		body		string	false	"Email"
//	@Param			phone		body		string	false	"Phone"
//	@Param			hire_date	body		string	false	"Hire date"
//	@Param			department_id	body		int		false	"Department ID"
//	@Param			can_login	body		string	false	"Can login"
//	@Param			password	body		string	false	"Password"
//	@Success		200			{object}	StatusOK
//	@Failure		400			{object}	StatusBadRequest
//	@Failure		401			{object}	StatusUnauthorized
//	@Security		apikey
//	@Router			/api/v1/employee/{id} [put]
func updateEmployee(c *gin.Context) {
	id := c.Param("id")
	var employee addEmployeeJson
	err := c.BindJSON(&employee)
	if err != nil {
		log.Printf("ERROR: Cannot bind JSON: %v", err)
		c.JSON(400, gin.H{
			"message": "Please provide a valid JSON",
		})
		return
	}

	// Check if id is a number
	if !isNumber(id) {
		handleError(c, 400, "ID must be a number", nil)
		return
	}

	key := ""
	value := ""

	// First name
	if employee.FirstName != "" {
		// Check if name contains only letters
		if err := checkName(employee.FirstName); err != nil {
			handleError(c, 400, err.Error(), err)
			return
		}
		key = "first_name"
		value = employee.FirstName
	}

	// Last name
	if employee.LastName != "" {
		// Check if name contains only letters
		if err := checkLastName(employee.LastName); err != nil {
			handleError(c, 400, err.Error(), err)
			return
		}
		key = "last_name"
		value = employee.LastName
	}

	// Email
	if employee.Email != "" {
		// Check for spaces in the email
		if err := checkEmail(employee.Email); err != nil {
			handleError(c, 400, err.Error(), err)
			return
		}
		key = "email"
		value = employee.Email
	}

	// Phone
	if employee.Phone != "" {
		// Check if the phone number contains only numbers
		if err := checkPhoneNumber(employee.Phone); err != nil {
			handleError(c, 400, err.Error(), err)
			return
		}
		key = "phone"
		value = employee.Phone
	}

	// Hire date
	if employee.HireDate != "" {
		key = "hire_date"
		value = employee.HireDate
	}

	// Department ID
	if employee.DepartmentID != 0 {
		// Check if the department exists
		if err := checkDepartmentExists(employee.DepartmentID); err != nil {
			handleError(c, 400, "Department does not exist", err)
			return
		}
		key = "department_id"
		value = fmt.Sprint(employee.DepartmentID)
	}

	// Can login is sent via json
	if employee.CanLogin != "" {
		key = "can_login"
		if employee.CanLogin == "true" {
			value = "1"
		} else {
			value = "0"
		}
	}

	// Check if the key is empty
	if value == "" {
		handleError(c, 400, "Please provide a value to update", nil)
		return
	}

	// Update in DB
	if updateEmployeeInDB(id, key, value) {
		c.JSON(200, gin.H{
			"message": "User updated",
			"id":      id,
			"key":     key,
			"value":   value,
		})
		return
	}

	handleError(c, 400, "User not updated", nil)
}

// does the actual database query
func updateEmployeeInDB(id string, key string, value string) bool {
	var err error

	// First name
	if key == "first_name" {
		// Update the first name in the database
		query := `UPDATE employee SET first_name = ?, edit_date = CURRENT_TIMESTAMP WHERE id = ?`
		_, err = db.Query(query, value, id)
		if err != nil {
			log.Printf("ERROR: Cannot update first name: %v", err)
			return false
		}
		return true
	}

	// Last name
	if key == "last_name" {
		// Update the last name in the database
		query := `UPDATE employee SET last_name = ?, edit_date = CURRENT_TIMESTAMP WHERE id = ?`
		_, err = db.Query(query, value, id)
		if err != nil {
			log.Printf("ERROR: Cannot update last name: %v", err)
			return false
		}
		return true
	}

	// Email
	if key == "email" {
		// Update the email in the database
		query := `UPDATE employee SET email = ?, edit_date = CURRENT_TIMESTAMP WHERE id = ?`
		_, err = db.Query(query, value, id)
		if err != nil {
			log.Printf("ERROR: Cannot update email: %v", err)
			return false
		}
		return true
	}

	// Phone
	if key == "phone" {
		// Update the phone in the database
		query := `UPDATE employee SET phone = ?, edit_date = CURRENT_TIMESTAMP WHERE id = ?`
		_, err = db.Query(query, value, id)
		if err != nil {
			log.Printf("ERROR: Cannot update phone: %v", err)
			return false
		}
		return true
	}

	// Hire date
	if key == "hire_date" {
		// Update the hire date in the database
		query := `UPDATE employee SET hire_date = ?, edit_date = CURRENT_TIMESTAMP WHERE id = ?`
		_, err = db.Query(query, value, id)
		if err != nil {
			log.Printf("ERROR: Cannot update hire date: %v", err)
			return false
		}
		return true
	}

	// Department ID
	if key == "department_id" {
		// Update the department ID in the database
		query := `UPDATE employee SET department_id = ?, edit_date = CURRENT_TIMESTAMP WHERE id = ?`
		_, err = db.Query(query, value, id)
		if err != nil {
			log.Printf("ERROR: Cannot update department ID: %v", err)
			return false
		}
		return true
	}

	// Can login
	if key == "can_login" {
		// Update the can login in the database
		query := `UPDATE employee SET can_login = ?, edit_date = CURRENT_TIMESTAMP WHERE id = ?`
		_, err = db.Query(query, value, id)
		if err != nil {
			log.Printf("ERROR: Cannot update can login: %v", err)
			return false
		}
		return true
	}

	return false
}

// DeleteEmployee
//
//	@Summary		Delete employee
//	@Description	Delete an employee from the database
//	@Tags			employees
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"Employee ID"
//	@Success		200			{object}	StatusOK
//	@Failure		400			{object}	StatusBadRequest
//	@Failure		401			{object}	StatusUnauthorized
//	@Security		apikey
//	@Router			/api/v1/employee/{id} [delete]
func deleteEmployee(c *gin.Context) {
	id := c.Param("id")
	employee := getEmployeeFromDB(id)

	// Check if id is a number
	if !isNumber(id) {
		handleError(c, 400, "ID must be a number", nil)
		return
	}

	if len(employee) == 0 {
		c.JSON(400, gin.H{
			"message": "User not found",
		})
		return
	}

	// Check if the employee is a manager
	departments := getDepartmentFromDB("")
	for _, dep := range departments {
		if dep.Manager == employee[0].FirstName+" "+employee[0].LastName {
			c.JSON(400, gin.H{
				"message":  "User is a manager",
				"employee": employee,
			})
			return
		}
	}

	if deleteEmployeeFromDB(id) {
		c.JSON(200, gin.H{
			"message":  "User deleted",
			"employee": employee,
		})
		return
	}

	c.JSON(400, gin.H{
		"message":  "User not deleted",
		"employee": employee,
	})
}

// does the actual database query
func deleteEmployeeFromDB(id string) bool {
	var err error

	// Delete the employee from the database
	query := `DELETE FROM employee WHERE id = ?`
	_, err = db.Query(query, id)
	if err != nil {
		log.Printf("ERROR: Cannot delete employee: %v", err)
		return false
	}

	return true
}

//
// Departments
//

// struct for the JSON body of the updateEmployee endpoint
type getDepartmentJson struct {
	ID       int    `json:"id" example:"1"`
	Name     string `json:"name" example:"Front Office"`
	Location string `json:"location" example:"Amsterdam"`
	Manager  string `json:"manager" example:"John Doe"`
}

// struct for the JSON body of the addDepartment endpoint
type addDepartmentJson struct {
	Name       string `json:"name" example:"Front Office"`
	LocationID int    `json:"location_id" example:"1"`
	ManagerID  int    `json:"manager_id" example:"1"`
}

// GetDepartment
//
//	@Summary		Show department(s)
//	@Description	Get all information from all departments or a specific department
//	@Tags			departments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	false	"Department ID"
//	@Success		200	{object}	getDepartmentJson
//	@Failure		401	{object}	StatusUnauthorized
//	@Security		apikey
//	@Router			/api/v1/department/{id} [get]
func getDepartment(c *gin.Context) {
	// If ID is provided, return the department with that ID
	id := c.Param("id")
	if id != "" {
		// Check if id is a number
		if !isNumber(id) {
			handleError(c, 400, "ID must be a number", nil)
			return
		}
		results := getDepartmentFromDB(id)
		c.JSON(http.StatusOK, results)
		return
	}

	results := getDepartmentFromDB("")
	c.JSON(http.StatusOK, results)
}

// does the actual database query
func getDepartmentFromDB(givenID string) []getDepartmentJson {
	rows := &sql.Rows{}
	departmentList := []getDepartmentJson{}
	var err error

	if givenID != "" {
		// Get the department with the provided ID
		// Inner join with location and manager
		query := `SELECT department.id, department.name, location.name as Location, CONCAT(manager.first_name, ' ', manager.last_name) as Manager FROM department
		INNER JOIN location ON department.location_id = location.id
		INNER JOIN employee AS manager ON department.manager_id = manager.id
		WHERE department.id = ?`

		rows, err = db.Query(query, givenID)
		if err != nil {
			log.Printf("ERROR: Cannot create table: %v", err)
		}
	} else {
		// Get all departments
		query := `SELECT department.id, department.name, location.name as Location, CONCAT(manager.first_name, ' ', manager.last_name) as Manager FROM department
		INNER JOIN location ON department.location_id = location.id
		INNER JOIN employee AS manager ON department.manager_id = manager.id
		ORDER BY department.id`

		rows, err = db.Query(query)
		if err != nil {
			log.Printf("ERROR: Cannot create table: %v", err)
		}
	}

	// Create a list of departments
	var name, location, manager string
	var id int

	for rows.Next() {
		err := rows.Scan(&id, &name, &location, &manager)
		if err != nil {
			log.Printf("ERROR: Cannot scan row: %v", err)
		}

		// Create a department object
		departmentinfo := getDepartmentJson{
			ID:       id,
			Name:     name,
			Location: location,
			Manager:  manager,
		}

		// Append the department object to the list
		departmentList = append(departmentList, departmentinfo)
	}

	// return as json
	return departmentList
}

// AddDepartment
//
//	@Summary		Add department
//	@Description	Add a new department to the database
//	@Tags			departments
//	@Accept			json
//	@Produce		json
//	@Param			name			body		string	true	"Name"
//	@Param			location_id		body		int		true	"Location ID"
//	@Param			manager_id		body		int		true	"Manager ID"
//	@Success		200				{object}	StatusOK
//	@Failure		400				{object}	StatusBadRequest
//	@Failure		401				{object}	StatusUnauthorized
//	@Security		apikey
//	@Router			/api/v1/department [post]
func addDepartment(c *gin.Context) {
	var department addDepartmentJson
	err := c.BindJSON(&department)
	if err != nil {
		log.Printf("ERROR: Cannot bind JSON: %v", err)
		c.JSON(400, gin.H{
			"message": "Please provide a valid JSON",
		})
		return
	}

	// Check if the required fields are provided
	if department.Name == "" || department.LocationID == 0 || department.ManagerID == 0 {
		c.JSON(400, gin.H{
			"message": "Please provide all required fields",
		})
		return
	}

	// Check if the location exists
	location := getLocationFromDB(fmt.Sprint(department.LocationID))
	if len(location) == 0 {
		c.JSON(400, gin.H{
			"message": "Location does not exist",
		})
		return
	}

	// Check if the manager exists
	manager := getEmployeeFromDB(fmt.Sprint(department.ManagerID))
	if len(manager) == 0 {
		c.JSON(400, gin.H{
			"message": "Manager does not exist",
		})
		return
	}

	// Add the department to the database and return the ID
	newID := addDepartmentToDB(department)
	c.JSON(200, gin.H{
		"message": "Department added",
		"id":      newID,
	})
}

// does the actual database query
func addDepartmentToDB(department addDepartmentJson) int {
	var err error

	// Insert the department into the database
	query := `INSERT INTO department (name, location_id, manager_id) VALUES (?, ?, ?)`
	_, err = db.Query(query, department.Name, department.LocationID, department.ManagerID)
	if err != nil {
		log.Printf("ERROR: Cannot insert department: %v", err)
	}

	// Get the ID of the inserted department
	query = `SELECT id FROM department WHERE name = ? AND location_id = ? AND manager_id = ?`
	rows, err := db.Query(query, department.Name, department.LocationID, department.ManagerID)
	if err != nil {
		log.Printf("ERROR: Cannot get ID: %v", err)
	}

	// Get the ID from the row
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Printf("ERROR: Cannot scan row: %v", err)
		}
	}

	// Send the ID back with the return statement
	return id
}

// UpdateDepartment
//
//	@Summary		Update department
//	@Description	Update a department in the database
//	@Tags			departments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"Department ID"
//	@Param			name	body		string	false	"Name"
//	@Param			location_id	body		int		false	"Location ID"
//	@Param			manager_id	body		int		false	"Manager ID"
//	@Success		200			{object}	StatusOK
//	@Failure		400			{object}	StatusBadRequest
//	@Failure		401			{object}	StatusUnauthorized
//	@Security		apikey
//	@Router			/api/v1/department/{id} [put]
func updateDepartment(c *gin.Context) {
	id := c.Param("id")

	// Check if id is a number
	if !isNumber(id) {
		handleError(c, 400, "ID must be a number", nil)
		return
	}

	c.JSON(200, gin.H{
		"message": "PUT /departments/:id",
		"id":      id,
	})
}

// DeleteDepartment
//
//	@Summary		Delete department
//	@Description	Delete a department from the database
//	@Tags			departments
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"Department ID"
//	@Success		200			{object}	StatusOK
//	@Failure		400			{object}	StatusBadRequest
//	@Failure		401			{object}	StatusUnauthorized
//	@Security		apikey
//	@Router			/api/v1/department/{id} [delete]
func deleteDepartment(c *gin.Context) {
	id := c.Param("id")

	// Check if id is a number
	if !isNumber(id) {
		handleError(c, 400, "ID must be a number", nil)
		return
	}

	department := getDepartmentFromDB(id)

	// Check if the department exists
	if len(department) == 0 {
		c.JSON(400, gin.H{
			"message": "Department not found",
		})
		return
	}

	// If check fails, return an error
	if !deleteDepartmentFromDB(id) {
		c.JSON(400, gin.H{
			"message":    "Department not deleted",
			"department": department,
		})
		return
	}

	// Check if the department has employees in it
	employees := getEmployeeFromDB("")
	for _, emp := range employees {
		if emp.Department == department[0].Name {
			c.JSON(400, gin.H{
				"message":    "Department has employees",
				"department": department,
			})
			return
		}
	}

	// If checks passes, return success
	c.JSON(200, gin.H{
		"message":    "Department deleted",
		"department": department,
	})
}

// does the actual database query
func deleteDepartmentFromDB(id string) bool {
	var err error

	// Delete the department from the database
	query := `DELETE FROM department WHERE id = ?`
	_, err = db.Query(query, id)
	if err != nil {
		log.Printf("ERROR: Cannot delete department: %v", err)
		return false
	}

	return true
}

//
// Location
//

// struct for the JSON body of the getLocation endpoint
type getLocationJson struct {
	ID         int    `json:"id" example:"1"`
	Name       string `json:"name" example:"Holiday Park"`
	Address    string `json:"address" example:"Street 1"`
	PostalCode string `json:"postal_code" example:"1234 AB"`
	City       string `json:"city" example:"Amsterdam"`
	Country    string `json:"country" example:"Netherlands"`
}

// struct for the JSON body of the addLocation endpoint
type addLocationJson struct {
	Name       string `json:"name" example:"Holiday Park"`
	Address    string `json:"address" example:"Street 1"`
	PostalCode string `json:"postal_code" example:"1234 AB"`
	City       string `json:"city" example:"Amsterdam"`
	Country    string `json:"country" example:"Netherlands"`
}

// GetLocation
//
//	@Summary		Show location(s)
//	@Description	Get all information from all locations or a specific location
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	false	"Location ID"
//	@Success		200	{object}	getLocationJson
//	@Failure		401	{object}	StatusUnauthorized
//	@Security		apikey
//	@Router			/api/v1/location/{id} [get]
func getLocation(c *gin.Context) {
	// If ID is provided, return the location with that ID
	id := c.Param("id")
	if id != "" {
		// Check if id is a number
		if !isNumber(id) {
			handleError(c, 400, "ID must be a number", nil)
			return
		}
		results := getLocationFromDB(id)
		c.JSON(http.StatusOK, results)
	}
	results := getLocationFromDB("")
	c.JSON(http.StatusOK, results)
}

// does the actual database query
func getLocationFromDB(givenID string) []getLocationJson {
	rows := &sql.Rows{}

	var err error

	if givenID != "" {
		// Get the location with the provided ID
		query := `SELECT * FROM location WHERE id = ?`
		rows, err = db.Query(query, givenID)
		if err != nil {
			log.Printf("ERROR: Cannot get individual location data: %v", err)
		}
	} else {
		// Get all locations
		query := `SELECT * FROM location`
		rows, err = db.Query(query)
		if err != nil {
			log.Printf("ERROR: Cannot get location data: %v", err)
		}
	}

	// Create a list of locations
	var name, address, postal_code, city, country string
	var id int
	locationList := []getLocationJson{}

	for rows.Next() {
		err := rows.Scan(&id, &name, &address, &postal_code, &city, &country)
		if err != nil {
			log.Printf("ERROR: Cannot scan row: %v", err)
		}

		// Create a location object
		locationinfo := getLocationJson{
			ID:         id,
			Name:       name,
			Address:    address,
			PostalCode: postal_code,
			City:       city,
			Country:    country,
		}

		// Append the location object to the json list
		locationList = append(locationList, locationinfo)
	}

	// return as json
	return locationList
}

// AddLocation
//
//	@Summary		Add location
//	@Description	Add a new location to the database
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			name			body		string	true	"Name"
//	@Param			address			body		string	true	"Address"
//	@Param			postal_code		body		string	true	"Postal code"
//	@Param			city			body		string	true	"City"
//	@Param			country			body		string	true	"Country"
//	@Success		200				{object}	StatusOK
//	@Failure		400				{object}	StatusBadRequest
//	@Failure		401				{object}	StatusUnauthorized
//	@Security		apikey
//	@Router			/api/v1/location [post]
func addLocation(c *gin.Context) {
	var location addLocationJson
	err := c.BindJSON(&location)
	if err != nil {
		log.Printf("ERROR: Cannot bind JSON: %v", err)
		c.JSON(400, gin.H{
			"message": "Please provide a valid JSON",
		})
		return
	}

	// Check if the required fields are provided
	if location.Name == "" || location.Address == "" || location.PostalCode == "" || location.City == "" || location.Country == "" {
		c.JSON(400, gin.H{
			"message": "Please provide all required fields",
		})
		return
	}

	// Check if the location already exists
	locations := getLocationFromDB("")
	for _, loc := range locations {
		if loc.Name == location.Name && loc.Address == location.Address && loc.PostalCode == location.PostalCode && loc.City == location.City && loc.Country == location.Country {
			c.JSON(400, gin.H{
				"message": "Location already exists",
			})
			return
		}
	}

	// Name must contain only letters
	if err := checkName(location.Name); err != nil {
		handleError(c, 400, err.Error(), err)
		return
	}

	// Address must contain only letters and numbers
	if err := checkAddress(location.Address); err != nil {
		handleError(c, 400, err.Error(), err)
		return
	}

	// City must contain only letters
	if err := checkName(location.City); err != nil {
		handleError(c, 400, err.Error(), err)
		return
	}

	// Add the location to the database and return the ID
	newID := addLocationToDB(location)
	c.JSON(200, gin.H{
		"message": "Location added",
		"id":      newID,
	})
}

// does the actual database query
func addLocationToDB(location addLocationJson) int {
	var err error

	// Insert the location into the database
	query := `INSERT INTO location (name, address, postal_code, city, country) VALUES (?, ?, ?, ?, ?)`
	_, err = db.Query(query, location.Name, location.Address, location.PostalCode, location.City, location.Country)
	if err != nil {
		log.Printf("ERROR: Cannot insert location: %v", err)
	}

	// Get the ID of the inserted location
	query = `SELECT id FROM location WHERE name = ? AND address = ? AND postal_code = ? AND city = ? AND country = ?`
	rows, err := db.Query(query, location.Name, location.Address, location.PostalCode, location.City, location.Country)
	if err != nil {
		log.Printf("ERROR: Cannot get ID: %v", err)
	}

	// Get the ID from the row
	var id int
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Printf("ERROR: Cannot scan row: %v", err)
		}
	}

	// Send the ID back with the return statement
	return id
}

// UpdateLocation
//
//	@Summary		Update location
//	@Description	Update a location in the database
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"Location ID"
//	@Param			name	body		string	false	"Name"
//	@Param			address	body		string	false	"Address"
//	@Param			postal_code	body		string	false	"Postal code"
//	@Param			city	body		string	false	"City"
//	@Param			country	body		string	false	"Country"
//	@Success		200			{object}	StatusOK
//	@Failure		400			{object}	StatusBadRequest
//	@Failure		401			{object}	StatusUnauthorized
//	@Security		apikey
//	@Router			/api/v1/location/{id} [put]
func updateLocation(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{
		"message": "PUT /location/:id",
		"id":      id,
	})
}

// DeleteLocation
//
//	@Summary		Delete location
//	@Description	Delete a location from the database
//	@Tags			locations
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"Location ID"
//	@Success		200			{object}	StatusOK
//	@Failure		400			{object}	StatusBadRequest
//	@Failure		401			{object}	StatusUnauthorized
//	@Security		apikey
//	@Router			/api/v1/location/{id} [delete]
func deleteLocation(c *gin.Context) {
	id := c.Param("id")
	location := getLocationFromDB(id)

	// Check if the location exists
	if len(location) == 0 {
		c.JSON(400, gin.H{
			"message": "Location not found",
		})
		return
	}

	// Check if the location is used in a department
	departments := getDepartmentFromDB("")
	for _, dep := range departments {
		if dep.Location == location[0].Name {
			c.JSON(400, gin.H{
				"message": "Location is used in a department",
			})
			return
		}
	}

	// If check fails, return an error
	if !deleteLocationFromDB(id) {
		c.JSON(400, gin.H{
			"message":  "Location not deleted",
			"location": location,
		})
		return
	}

	// If checks passes, return success
	c.JSON(200, gin.H{
		"message": "Location is deleted",
		"id":      id,
	})
}

// does the actual database query
func deleteLocationFromDB(id string) bool {
	var err error

	// Delete the location from the database
	query := `DELETE FROM location WHERE id = ?`
	_, err = db.Query(query, id)
	if err != nil {
		log.Printf("ERROR: Cannot delete location: %v", err)
		return false
	}
	return true
}
