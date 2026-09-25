module github.com/sokosplit/sokosplit-core-api

go 1.22

require (
	github.com/gin-gonic/gin v1.10.0
	gorm.io/driver/postgres v1.5.9
	gorm.io/gorm v1.25.10
	github.com/google/uuid v1.6.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/glebarez/sqlite v1.11.0 // test-only: pure-Go sqlite driver for handler tests
)
