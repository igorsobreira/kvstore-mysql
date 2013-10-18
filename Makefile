
test:
	@mysql -u root -e "DROP DATABASE IF EXISTS kvstore_mysql_test" 1>/dev/null 2>&1 || (echo "Start mysql accepting connections from root without password" && exit 1)
	@mysql -u root -e "CREATE DATABASE kvstore_mysql_test"
	@mysql -u root -e "CREATE DATABASE IF NOT EXISTS kvstore_mysql_example"
	@go test -race -i
	@go test -race -v

lint:
	@golint `find . -name "*.go"`

fmt:
	@go fmt ./...
