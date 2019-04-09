THRIFT=thrift
THRIFT_GO=$(THRIFT) --gen go

gen-go/tutorial/tutorial.go: tutorial.thrift
	$(THRIFT_GO) $<
