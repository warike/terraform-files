# Contributing to tfinit

First off, thank you for considering contributing to `tfinit`. It's people like you that make `tfinit` such a great tool.

Following these guidelines helps to communicate that you respect the time of the developers managing and developing this open source project. In return, they should reciprocate that respect in addressing your issue, assessing changes, and helping you finalize your pull requests.

## Getting Started

To get started with development, you'll need a working Go environment.

1.  **Fork and Clone the repository**
    ```bash
    git clone https://github.com/warike/terraform-files.git
    cd terraform-files
    ```

2.  **Install Dependencies**
    This project uses Go modules. Dependencies will be downloaded automatically when you build or test the project.
    ```bash
    go mod tidy
    ```

3.  **Run Tests**
    We use `make` to simplify testing. To run all checks, including unit and integration tests, run:
    ```bash
    make test
    ```

4.  **Build the binary**
    To compile the `tfinit` binary for your local system:
    ```bash
    make build
    ```
    The binary will be available at `build/tfinit`.

## Pull Request Process

1.  Ensure any install or build dependencies are removed before the end of the layer when doing a build.
2.  Update the README.md with details of changes to the interface, this includes new environment variables, exposed ports, useful file locations and container parameters.
3.  Increase the version numbers in any examples files and the README.md to the new version that this Pull Request would represent. The versioning scheme we use is [SemVer](http://semver.org/).
4.  You may merge the Pull Request in once you have the sign-off of two other developers, or if you do not have permission to do that, you may request the second reviewer to merge it for you.

## Code of Conduct

This project and everyone participating in it is governed by the [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior.
