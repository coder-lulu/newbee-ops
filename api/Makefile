SERVICE := newbee-ops-api
PROJECT_STYLE := go_zero

.PHONY: gen-api

gen-api:
	@echo "Generating API server from desc/all.api ..."
	goctl api go -api desc/all.api -dir . --style=$(PROJECT_STYLE) --home ~/.goctl
	@echo "✅ API generation complete. Check internal/handler/routes.go for updated routes."

