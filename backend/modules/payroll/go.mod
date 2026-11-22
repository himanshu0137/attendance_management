module hunrx.com/attendance_management/modules/payroll

go 1.23

require (
	hunrx.com/attendance_management v0.0.0
	github.com/caarlos0/env/v11 v11.4.2
	github.com/gin-gonic/gin v1.11.0
	github.com/jackc/pgx/v5 v5.13.0
	github.com/pressly/goose/v3 v3.30.1
)

replace hunrx.com/attendance_management => ../..
