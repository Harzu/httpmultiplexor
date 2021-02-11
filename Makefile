REPOSITORY := "fanyshu"
APP_NAME := "http-multiplexor"
VERSION := $(if $(TAG),$(TAG),$(if $(BRANCH_NAME),$(BRANCH_NAME),$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)))

build:
	@docker build -t $(REPOSITORY)/$(APP_NAME):$(VERSION) .

start:
	@go run ./cmd

check:
	@curl -X POST http://localhost:8080/pages -H 'Content-Type: application/json' \
		-d '["https://google.com", "https://github.com", "https://gitlab.com", "https://yandex.ru", "https://mail.ru", "https://google.com", "https://github.com", "https://gitlab.com", "https://yandex.ru", "https://mail.ru"]'

check_with_error:
	@curl -X POST http://localhost:8080/pages -H 'Content-Type: application/json' \
    		-d '["https://google.com", "https://github.com", "https://gitlab.com", "http://localhost:90099/", "https://yandex.ru", "https://mail.ru", "https://google.com", "https://github.com", "https://gitlab.com", "https://yandex.ru", "https://mail.ru"]'

rate_limit_test_success:
	@vegeta attack -duration=10s -rate=5 -targets=vegeta.conf | vegeta report

rate_limit_test_exceed_limit:
	@vegeta attack -duration=10s -rate=85 -workers=1 -targets=vegeta.conf | vegeta report