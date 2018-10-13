    # Go parameters
    GOCMD=go
    GLIDE=glide
		CMDDIR=./cmd
    OUTDIR=./out
		VENDOR=./vendor
    BINARY_NAME=access-stats
    GOBUILD=$(GOCMD) build
    GOCLEAN=$(GOCMD) clean
    GOTEST=$(GOCMD) test
    GOGET=$(GOCMD) get
    
    all: test build
    build: 
			$(GOBUILD) -o $(OUTDIR)/$(BINARY_NAME) -v $(CMDDIR)
    test: 
			$(GOTEST) -v ./...
    clean: 
			$(GOCLEAN)
			rm -f $(BINARY_NAME)
			rm -f $(VENDOR)
    run: build	
			$(OUTDIR)/$(BINARY_NAME) -src=$(src) -alertThreshold=$(alertThreshold)
    deps:
			$(GOGET) github.com/Masterminds/glide
			$(GLIDE) i