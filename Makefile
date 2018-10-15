    # Go parameters
    GOCMD=go
    GLIDE=glide
		CMDDIR=./cmd
    OUTDIR=./out
		VENDOR=./vendor
    SCRIPTSDIR=./scripts
    BINARY_NAME=access-stats
    GOBUILD=$(GOCMD) build
    GOCLEAN=$(GOCMD) clean
    GOTEST=$(GOCMD) test
    GOGET=$(GOCMD) get
    GORUN=$(GOCMD) run
    
    all: test build
    build: deps
			$(GOBUILD) -o $(OUTDIR)/$(BINARY_NAME) -v $(CMDDIR)
    build-linux: deps
			GOOS=linux $(GOBUILD) -o $(OUTDIR)/$(BINARY_NAME) -v $(CMDDIR)
    test: 
			$(GOTEST) -v ./...
    clean: 
			$(GOCLEAN)
			rm -rf $(OUTDIR)
			rm -rf $(VENDOR)
    run: build	
			$(OUTDIR)/$(BINARY_NAME) -src=$(src) -alertThreshold=$(alertThreshold)
    deps:
			$(GOGET) github.com/Masterminds/glide
			$(GLIDE) i
    sample-data:
			$(GORUN) $(SCRIPTSDIR)/alert-recover.go -src=$(src)