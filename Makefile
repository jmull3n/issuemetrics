test:
	@go test -v ./...

run:
	@echo "running on port 8000"
	@docker run -p 8000:8000 jmull3n/issuemetrics serve