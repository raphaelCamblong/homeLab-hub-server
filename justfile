set shell := ["bash", "-c"]

APP_NAME := "homelab_server"
main_path := "homelab.com/homelab-server/homeLab-server/cmd"

build_args := "-v"

build:
    @echo "Building the application {{APP_NAME}}..."
    go build {{build_args}} -o {{APP_NAME}} {{main_path}}

# Run the application
run: build
    @echo "Running the application..."
    ./{{APP_NAME}}

migrate:
    @echo "Not implemented yet..."

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
