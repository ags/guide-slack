BINS = $(shell ls cmd)

.PHONY: all
all: $(BINS)

.PHONY: $(BINS)
$(BINS):
	GOOS=linux GOARCH=amd64 go build -o tmp/$@ cmd/$@/main.go
	zip -j -o tmp/$@.zip tmp/$@

.PHONY:
clean:
	rm -rf tmp
