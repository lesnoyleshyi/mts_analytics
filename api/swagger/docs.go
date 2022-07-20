// Package swagger GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package swagger

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/agreed": {
            "get": {
                "description": "Get count of all signed tasks.",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "count"
                ],
                "summary": "get signed tasks count",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/canceled": {
            "get": {
                "description": "Get count of all tasks which are not signed (rejected and \"in process\").",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "count"
                ],
                "summary": "get unsigned tasks count",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/total_time": {
            "get": {
                "description": "Get total signition time in seconds of particular task by its id",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "time"
                ],
                "summary": "get signition time of task",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "uuid of task",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "example: 'time in sec: 100500'",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "returns in case user input is invalid",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "returns in case server can't retrieve signition time",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "": {
            "description": "Cookie authorization requires a pair of Access and Refresh tokens. Get them in auth service. Url below is the place where you can do it.",
            "type": "oauth2",
            "flow": "application",
            "tokenUrl": "auth-service/v1/auth/validate"
        }
    },
    "x-extension-openapi": {
        "example": "value on a json format"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.0.1",
	Host:             "localhost:80",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Analytics service",
	Description:      "Analytics service knows lifecycle of every task.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
