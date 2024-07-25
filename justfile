APP_NAME := homelab_server

set shell := ["bash", "-c"]

build:
    @echo "Building the application..."
    go build -o {{APP_NAME}} .

# Run the application
run: build
    @echo "Running the application..."
    ./{{APP_NAME}}

# Test the application
test:
    @echo "Running tests..."
    go test ./...

# Clean the build directory
clean:
    @echo "Cleaning build directory..."
    go clean
    rm -rf {{APP_NAME}}

# Format the code
fmt:
    @echo "Formatting the code..."
    go fmt ./...

# Download dependencies
deps:
    @echo "Downloading dependencies..."
    go get -u ./...

# Tidy up the go.mod file
tidy:
    @echo "Tidying up the go.mod file..."
    go mod tidy
