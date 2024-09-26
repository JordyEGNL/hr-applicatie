// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Jordy Hoebergen",
            "email": "j.hoebergen@student.fontys.nl"
        },
        "license": {
            "name": "GPL-3.0",
            "url": "https://www.gnu.org/licenses/gpl-3.0.en.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/department": {
            "post": {
                "security": [
                    {
                        "apikey": []
                    }
                ],
                "description": "Add a new department to the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "departments"
                ],
                "summary": "Add department",
                "parameters": [
                    {
                        "description": "Name",
                        "name": "name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Location ID",
                        "name": "location_id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "Manager ID",
                        "name": "manager_id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.StatusOK"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.StatusBadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.StatusUnauthorized"
                        }
                    }
                }
            }
        },
        "/api/v1/department/{id}": {
            "get": {
                "security": [
                    {
                        "apikey": []
                    }
                ],
                "description": "Get all information from all departments or a specific department",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "departments"
                ],
                "summary": "Show department(s)",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Department ID",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.getDepartmentJson"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.StatusUnauthorized"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "apikey": []
                    }
                ],
                "description": "Update a department in the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "departments"
                ],
                "summary": "Update department",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Department ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Name",
                        "name": "name",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Location ID",
                        "name": "location_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "Manager ID",
                        "name": "manager_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.StatusOK"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.StatusBadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.StatusUnauthorized"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "apikey": []
                    }
                ],
                "description": "Delete a department from the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "departments"
                ],
                "summary": "Delete department",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Department ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.StatusOK"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.StatusBadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.StatusUnauthorized"
                        }
                    }
                }
            }
        },
        "/api/v1/employee": {
            "post": {
                "security": [
                    {
                        "apikey": []
                    }
                ],
                "description": "Add a new employee to the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "employees"
                ],
                "summary": "Add employee",
                "parameters": [
                    {
                        "description": "First name",
                        "name": "first_name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Last name",
                        "name": "last_name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Email",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Phone",
                        "name": "phone",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Hire date",
                        "name": "hire_date",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Department ID",
                        "name": "department_id",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "Can login",
                        "name": "can_login",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.StatusOK"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.StatusBadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.StatusUnauthorized"
                        }
                    }
                }
            }
        },
        "/api/v1/employee/{id}": {
            "get": {
                "security": [
                    {
                        "apikey": []
                    }
                ],
                "description": "Get all information from all employees or a specific employee",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "employees"
                ],
                "summary": "Show employee(s)",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Employee ID",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.getEmployeeJson"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.StatusUnauthorized"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "apikey": []
                    }
                ],
                "description": "Update an employee in the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "employees"
                ],
                "summary": "Update employee",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Employee ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "First name",
                        "name": "first_name",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Last name",
                        "name": "last_name",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Email",
                        "name": "email",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Phone",
                        "name": "phone",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Hire date",
                        "name": "hire_date",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Department ID",
                        "name": "department_id",
                        "in": "body",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    {
                        "description": "Can login",
                        "name": "can_login",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Password",
                        "name": "password",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.StatusOK"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.StatusBadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.StatusUnauthorized"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "apikey": []
                    }
                ],
                "description": "Delete an employee from the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "employees"
                ],
                "summary": "Delete employee",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Employee ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.StatusOK"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.StatusBadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.StatusUnauthorized"
                        }
                    }
                }
            }
        },
        "/api/v1/location": {
            "post": {
                "security": [
                    {
                        "apikey": []
                    }
                ],
                "description": "Add a new location to the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "locations"
                ],
                "summary": "Add location",
                "parameters": [
                    {
                        "description": "Name",
                        "name": "name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Address",
                        "name": "address",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Postal code",
                        "name": "postal_code",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "City",
                        "name": "city",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Country",
                        "name": "country",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.StatusOK"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.StatusBadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.StatusUnauthorized"
                        }
                    }
                }
            }
        },
        "/api/v1/location/{id}": {
            "get": {
                "security": [
                    {
                        "apikey": []
                    }
                ],
                "description": "Get all information from all locations or a specific location",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "locations"
                ],
                "summary": "Show location(s)",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Location ID",
                        "name": "id",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.getLocationJson"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.StatusUnauthorized"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "apikey": []
                    }
                ],
                "description": "Update a location in the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "locations"
                ],
                "summary": "Update location",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Location ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Name",
                        "name": "name",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Address",
                        "name": "address",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Postal code",
                        "name": "postal_code",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "City",
                        "name": "city",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "Country",
                        "name": "country",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.StatusOK"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.StatusBadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.StatusUnauthorized"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "apikey": []
                    }
                ],
                "description": "Delete a location from the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "locations"
                ],
                "summary": "Delete location",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Location ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.StatusOK"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.StatusBadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.StatusUnauthorized"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Autenhticate the user with the provided credentials to get access to the API and frontend",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Authenticate user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Username",
                        "name": "username",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "Password",
                        "name": "password",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.StatusOK"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.StatusUnauthorized"
                        }
                    }
                }
            }
        },
        "/logout": {
            "get": {
                "description": "Logout the user and remove the session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Logout user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.StatusOK"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/main.StatusBadRequest"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/main.StatusUnauthorized"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/main.StatusInternalServerError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.StatusBadRequest": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "main.StatusInternalServerError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                }
            }
        },
        "main.StatusOK": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "main.StatusUnauthorized": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "main.getDepartmentJson": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "location": {
                    "type": "string",
                    "example": "Amsterdam"
                },
                "manager": {
                    "type": "string",
                    "example": "John Doe"
                },
                "name": {
                    "type": "string",
                    "example": "Front Office"
                }
            }
        },
        "main.getEmployeeJson": {
            "type": "object",
            "properties": {
                "can_login": {
                    "type": "boolean",
                    "example": true
                },
                "creation_date": {
                    "type": "string",
                    "example": "2021-01-01 9:00:00"
                },
                "department": {
                    "type": "string",
                    "example": "1"
                },
                "edit_date": {
                    "type": "string",
                    "example": "2021-01-01 12:00:00"
                },
                "email": {
                    "type": "string",
                    "example": "j.doe@holiday-parks.eu"
                },
                "first_name": {
                    "type": "string",
                    "example": "John"
                },
                "hire_date": {
                    "type": "string",
                    "example": "2021-01-01"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "last_login": {
                    "type": "string",
                    "example": "2021-01-01 10:00:00"
                },
                "last_name": {
                    "type": "string",
                    "example": "Doe"
                },
                "location": {
                    "$ref": "#/definitions/main.getLocationJson"
                },
                "manager_email": {
                    "type": "string",
                    "example": "m.anager@holiday-parks.eu"
                },
                "phone": {
                    "type": "string",
                    "example": "06-12345678"
                }
            }
        },
        "main.getLocationJson": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string",
                    "example": "Street 1"
                },
                "city": {
                    "type": "string",
                    "example": "Amsterdam"
                },
                "country": {
                    "type": "string",
                    "example": "Netherlands"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "type": "string",
                    "example": "Holiday Park"
                },
                "postal_code": {
                    "type": "string",
                    "example": "1234 AB"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "127.0.0.1:5000",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "HR API",
	Description:      "This is the HR API for the Holiday Parks company.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
