#!/usr/bin/env sh
# Copyright 2024-present DevControl contributors.

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m'

# Function to check Linux distribution
check_distro() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        DISTRO=$ID
    elif [ -f /etc/lsb-release ]; then
        . /etc/lsb-release
        DISTRO=$DISTRIB_ID
    elif [ -f /etc/redhat-release ]; then
        DISTRO=$(cat /etc/redhat-release | cut -d ' ' -f 1)
    else
        DISTRO="Unknown"
    fi
    echo -e "\n${PURPLE}Linux Distribution: $DISTRO${NC}"
}

# Function to check Go and Docker installations
check_installations() {
    if ! command -v docker >/dev/null 2>&1; then
        echo -e "${RED}Docker is not installed.${NC}"
        echo "Please install Docker from https://docs.docker.com/get-docker/"
        INSTALL_DOCKER=true
    else
        echo -e "${GREEN}Docker is installed.${NC}"
    fi

    if ! docker buildx version >/dev/null 2>&1; then
        echo -e "${RED}Docker Buildx is not installed or not configured.${NC}"
        echo "Please install Docker Buildx from https://docs.docker.com/buildx/working-with-buildx/"
        INSTALL_BUILDX=true
    else
        echo -e "${GREEN}Docker Buildx is installed.${NC}"
    fi

    if ! command -v go >/dev/null 2>&1; then
        echo -e "${RED}Go is not installed.${NC}"
        echo "Please install Go from https://golang.org/dl/"
        INSTALL_GO=true
    else
        echo -e "${GREEN}Go is installed.${NC}"
    fi

    if [ "$INSTALL_DOCKER" = true ] || [ "$INSTALL_BUILDX" = true ] || [ "$INSTALL_GO" = true ]; then
        echo -e "\n${YELLOW}Please install the missing dependencies and run the script again.${NC}"
        exit 1
    fi
}

# Function to clone repository and set it up
setup_repository() {
    echo -e "\n${CYAN}Cloning repository...${NC}"
    REPO_DIR="$HOME/.devbox"
    git clone https://github.com/harshau007/devbox.git "$REPO_DIR" || {
        echo -e "${RED}Failed to clone repository.${NC}"
        exit 1
    }
    cd "$REPO_DIR" || exit 1

    echo -e "\n${CYAN}Setting up the project...${NC}"
    go mod tidy || {
        echo -e "${RED}Failed to set up the project.${NC}"
        remove_repo
        exit 1
    }

    echo -e "\n${BLUE}Building the project...${NC}"
    go build -o devctl || {
        echo -e "${RED}Failed to build the project.${NC}"
        remove_repo
        exit 1
    }

    echo -e "\n${BLUE}Copying binaries and files...${NC}"
    sudo cp devctl portdevctl startdevctl /usr/bin/ || {
        echo -e "${RED}Failed to copy binaries.${NC}"
        remove_repo
        exit 1
    }
    sudo mkdir -p /usr/local/share/devbox/ || {
        echo -e "${RED}Failed to create directory for config files.${NC}"
        remove_repo
        exit 1
    }
    sudo cp dockerfile settings.json setup.sh /usr/local/share/devbox/ || {
        echo -e "${RED}Failed to copy config files.${NC}"
        remove_repo
        exit 1
    }

    echo -e "\n${GREEN}Project setup completed!${NC}"
    echo -e "${YELLOW}Run 'devctl -h' for further details.${NC}"
}

# Function to remove the cloned repository
remove_repo() {
    rm -rf "$REPO_DIR"
}

# Function to handle script exit
trap remove_repo EXIT

check_distro
check_installations
setup_repository