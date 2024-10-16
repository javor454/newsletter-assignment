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
            "email": "javornicky.jiri@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/health/liveness": {
            "get": {
                "tags": [
                    "health"
                ],
                "summary": "Determines if app is running",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/api/health/readiness": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "health"
                ],
                "summary": "Determines if app is ready to receive traffic",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.HealthStatus"
                        }
                    }
                }
            }
        },
        "/api/v1/newsletters": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "newsletter"
                ],
                "summary": "Retrieve newsletter by creator's user ID",
                "parameters": [
                    {
                        "type": "string",
                        "default": "application/json",
                        "description": "application/json",
                        "name": "Content-Type",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "Bearer",
                        "description": "Bearer \u003ctoken\u003e",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "default": 10,
                        "description": "Number of items on page",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "default": 1,
                        "description": "Page number",
                        "name": "page_number",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully retrieved newsletters by user ID",
                        "schema": {
                            "$ref": "#/definitions/response.InternalNewsletter"
                        }
                    },
                    "400": {
                        "description": "Invalid request with detail",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Unexpected exception"
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "newsletter"
                ],
                "summary": "Create new newsletter",
                "parameters": [
                    {
                        "type": "string",
                        "default": "Bearer",
                        "description": "Bearer \u003ctoken\u003e",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Newsletter data to create",
                        "name": "Newsletter",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.CreateNewsletterRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Newsletter was successfully created"
                    },
                    "400": {
                        "description": "Invalid request with detail",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "500": {
                        "description": "Unexpected exception"
                    }
                }
            }
        },
        "/api/v1/newsletters/{newsletter_public_id}/subscriptions": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "public subscription"
                ],
                "summary": "Used to subscribe to newsletter by email",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Public newsletter identifier",
                        "name": "newsletter_public_id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Subscriber email address",
                        "name": "email",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.SubscribeToNewsletter"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Successfully subscribed to newsletter"
                    },
                    "400": {
                        "description": "Invalid request with detail",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "404": {
                        "description": "Newsletter not found",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "409": {
                        "description": "Already subscribed to newsletter",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Unexpected exception"
                    }
                }
            }
        },
        "/api/v1/newsletters/{public_id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "public newsletter"
                ],
                "summary": "Retrieve newsletter by its public ID",
                "parameters": [
                    {
                        "type": "string",
                        "default": "application/json",
                        "description": "application/json",
                        "name": "Content-Type",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Newsletter public ID",
                        "name": "public_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully retrieved newsletter by public ID",
                        "schema": {
                            "$ref": "#/definitions/response.PublicNewsletter"
                        }
                    },
                    "400": {
                        "description": "Invalid request with detail",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Unexpected exception"
                    }
                }
            }
        },
        "/api/v1/subscriptions/{email}/newsletters": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "public subscription"
                ],
                "summary": "Retrieve newsletter by subscriber's email",
                "parameters": [
                    {
                        "type": "string",
                        "default": "application/json",
                        "description": "application/json",
                        "name": "Content-Type",
                        "in": "header",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "default": 10,
                        "description": "Number of items on page",
                        "name": "page_size",
                        "in": "query",
                        "required": true
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "default": 1,
                        "description": "Page number",
                        "name": "page_number",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "test@test.com",
                        "description": "Subscribers email",
                        "name": "email",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully retrieved newsletters by subscriber email"
                    },
                    "400": {
                        "description": "Invalid request with detail",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Unexpected exception"
                    }
                }
            }
        },
        "/api/v1/unsubscribe": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "public subscription"
                ],
                "summary": "Used to unsubscribe from newsletter by email",
                "parameters": [
                    {
                        "type": "string",
                        "default": "application/json",
                        "description": "application/json",
                        "name": "Content-Type",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Public newsletter identifier",
                        "name": "newsletter_public_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Token to associate with subscription",
                        "name": "token",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully unsubscribed from newsletter"
                    },
                    "400": {
                        "description": "Invalid request with detail",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Unexpected exception"
                    }
                }
            }
        },
        "/api/v1/users/login": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "public user"
                ],
                "summary": "Login user, returning token for authorization",
                "parameters": [
                    {
                        "description": "Data for user login",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UserRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User successfully logged in"
                    },
                    "400": {
                        "description": "Invalid request with detail",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "401": {
                        "description": "Invalid credentials",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Unexpected exception"
                    }
                }
            }
        },
        "/api/v1/users/register": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "public user"
                ],
                "summary": "Register user, returning token for authorization",
                "parameters": [
                    {
                        "description": "Data for registering user",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UserRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User was successfully registered"
                    },
                    "400": {
                        "description": "Invalid request with detail",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "409": {
                        "description": "Email taken",
                        "schema": {
                            "$ref": "#/definitions/response.Error"
                        }
                    },
                    "500": {
                        "description": "Unexpected exception"
                    }
                }
            }
        }
    },
    "definitions": {
        "request.CreateNewsletterRequest": {
            "type": "object",
            "required": [
                "name"
            ],
            "properties": {
                "description": {
                    "type": "string",
                    "example": "Amazing news from the TikTok world. You would not believe number 4."
                },
                "name": {
                    "type": "string",
                    "example": "Tiktok News 420"
                }
            }
        },
        "request.SubscribeToNewsletter": {
            "type": "object",
            "required": [
                "email"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "test@test.com"
                }
            }
        },
        "request.UserRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "test@test.com"
                },
                "password": {
                    "type": "string",
                    "example": "Pa$$W0rD"
                }
            }
        },
        "response.Error": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Error description"
                }
            }
        },
        "response.HealthStatus": {
            "type": "object",
            "properties": {
                "indicators": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/response.Indicator"
                    }
                },
                "status": {
                    "type": "string",
                    "example": "healthy"
                }
            }
        },
        "response.Indicator": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "example": "postgres"
                },
                "status": {
                    "type": "string",
                    "example": "healthy"
                }
            }
        },
        "response.InternalNewsletter": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2024-09-20T23:16:32Z"
                },
                "description": {
                    "type": "string",
                    "example": "Some descriptive description"
                },
                "id": {
                    "type": "string",
                    "example": "1541c9c1-e43e-4527-850a-77f4e5be9599"
                },
                "name": {
                    "type": "string",
                    "example": "Newsletter name"
                },
                "public_id": {
                    "type": "string",
                    "example": "90c0a606-4429-44cc-9531-6f9cd038620a"
                }
            }
        },
        "response.PublicNewsletter": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string",
                    "example": "2024-09-20T23:16:32Z"
                },
                "description": {
                    "type": "string",
                    "example": "Some descriptive description"
                },
                "name": {
                    "type": "string",
                    "example": "Newsletter name"
                },
                "public_id": {
                    "type": "string",
                    "example": "90c0a606-4429-44cc-9531-6f9cd038620a"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Newsletter assignment",
	Description:      "Newsletter assignment for STRV.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
