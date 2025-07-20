## Define the project name
`APP_NAME="wk-terraform"`

## Create a 'build' folder to store the compiled binaries
`mkdir -p build`

## Compiling for macOS (Apple Silicon - ARM64)
`GOOS=darwin GOARCH=arm64 go build -o "build/${APP_NAME}_darwin_arm64" .`

## Compiling for macOS (Intel - x86_64)
`GOOS=darwin GOARCH=amd64 go build -o "build/${APP_NAME}_darwin_amd64" .`

## List the created binaries
`ls -lh build/`