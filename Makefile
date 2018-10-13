    # Go parameters
    GOCMD=go
    GLIDE=glide
		CMDDIR=./cmd
    OUTDIR=./out
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
    run: build	
			$(OUTDIR)/$(BINARY_NAME) -src=$(src)
    deps:
			$(GOGET) github.com/Masterminds/glide
			$(GLIDE) i