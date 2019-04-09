.DEFAULT_GOAL := demo

THRIFT=thrift
THRIFT_GO=$(THRIFT) --gen go
GO=go

gen-go/tutorial/tutorial.go: tutorial.thrift
	$(THRIFT_GO) $<

.PHONY: demo

demo: gen-go/tutorial/tutorial.go *.go
	$(GO) build
