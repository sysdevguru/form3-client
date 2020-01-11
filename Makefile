.PHONY: deps unit-test integration-test
deps:
	go get \
		github.com/pariz/gountries \
		golang.org/x/text/currency \
		gopkg.in/go-playground/validator.v9 \
		github.com/stretchr/testify/assert 

unit-test:
	cd ./form3go && go test --cover -v

integration-test:
	cd ./form3go && go test --tags=integration -v