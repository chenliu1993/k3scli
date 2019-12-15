CLIBIN := ./main

#$(CLIBIN):
#	@go build -o $(@)

.PHONY: build
build: echo $(CLIBIN)

$(CLIBIN):
	@go build -o k3scli
.PHONY: echo
echo: 
	@echo "-----------------------"
	@echo "beginning build..."
	@echo "-----------------------"

