package adapters

// @title Analytics service
// @version 0.0.1
// @description Analytics service knows lifecycle of every task.

// TODO this should be parsed from config
// @host localhost:80
// @BasePath /
// @query.collection.format multi
// @schemes http

// @securityDefinitions.oauth2.application
// @in header
// @name cookie authorization
//nolint:lll
// @description Cookie authorization requires a pair of Access and Refresh tokens. Get them in auth service. Url below is the place where you can do it.
// @tokenUrl auth-service/v1/auth/validate

// @x-extension-openapi {"example": "value on a json format"}
