#!/bin/bash
set -euo pipefail

# Warna untuk output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}  Employee Service - Build Script       ${NC}"
echo -e "${YELLOW}========================================${NC}"

# Direktori project (relative dari script)
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_DIR"

echo -e "${GREEN}[1/5] Checking Go version...${NC}"
go version

echo -e "${GREEN}[2/5] Downloading dependencies...${NC}"
go mod download

echo -e "${GREEN}[3/5] Running go vet (static analysis)...${NC}"
go vet ./...

echo -e "${GREEN}[4/5] Building application...${NC}"
go build -o bin/employee-server ./cmd/server/

echo -e "${GREEN}[5/5] Build complete!${NC}"
echo -e "${GREEN}Binary: bin/employee-server${NC}"

echo -e "\n${YELLOW}========================================${NC}"
echo -e "${YELLOW}  Running server:                        ${NC}"
echo -e "${YELLOW}  ./bin/employee-server                  ${NC}"
echo -e "${YELLOW}========================================${NC}"