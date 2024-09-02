set shell := ["bash", "-c"]

APP_NAME := "homelab_server"
main_path := "homelab.com/homelab-server/homeLab-server/cmd"
migration_sql_path := "homelab.com/homelab-server/database/migration"

build_args := "-v"

build:
    @echo "Building the application {{APP_NAME}}..."
    go build {{build_args}} -o {{APP_NAME}} {{main_path}}

# Run the application
run: build
    @echo "Running the application..."
    ./{{APP_NAME}}

migrate:
    @echo "Running migration..."
    go run {{migration_sql_path}}

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

# Docker

docker_user := "velocipastor"
# Build the docker image
docker-build:
    @echo "Building the docker image..."
    docker build -t {{APP_NAME}} .

docker-run:
    @echo "Running the docker image..."
    docker run -p 6000:6000 {{docker_user}}/{{APP_NAME}}:latest

docker-publish:
    docker tag {{APP_NAME}}:latest {{docker_user}}/{{APP_NAME}}:latest
    docker push {{docker_user}}/{{APP_NAME}}:latest

helm-install:
     cd cicd && helm install home-lab-backend ./homeLab-backend -f ./homeLab-backend/values.yaml -n home-lab

helm-delete:
     cd cicd && helm delete home-lab-backend -n home-lab
