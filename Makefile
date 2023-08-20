
DATABASE_DSN="postgres://anton:!anton321@localhost:5444/mart?sslmode=disable"

build:
	go build -buildvcs=false  -o ./cmd/gophermart/gophermart ./cmd/gophermart/main.go

tests:
	gophermarttest \
                -test.v -test.run=^TestGophermart$$ \
                -gophermart-binary-path=./cmd/gophermart/gophermart \
                -gophermart-host=localhost \
                -gophermart-port=8080 \
                -gophermart-database-uri=$(DATABASE_DSN) \
                -accrual-binary-path=./cmd/accrual/accrual \
                -accrual-host=localhost \
                -accrual-port=8088 \
                -accrual-database-uri=$(DATABASE_DSN)
