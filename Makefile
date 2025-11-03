PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := .

.PHONY: generate-proto
generate-proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_SRC)

# --- Configuration ---
BIN_DIR := bin
# Add .exe extension for Windows binaries
BIN_EXT := .exe

# Service Paths
SERVICE_A_PATH := services/api-gateway
SERVICE_B_PATH := services/user-service
SERVICE_C_PATH := services/chat-service
SERVICE_D_PATH := services/event-service
SERVICE_E_PATH := services/notification-service

# Service Names (Binaries will be named after these)
SERVICE_A_NAME := api-gateway
SERVICE_B_NAME := user-services
SERVICE_C_NAME := chat-service
SERVICE_D_NAME := event-service
SERVICE_E_NAME := notification-service

# Binary Paths (now include .exe)
BIN_A := $(BIN_DIR)/$(SERVICE_A_NAME)$(BIN_EXT)
BIN_B := $(BIN_DIR)/$(SERVICE_B_NAME)$(BIN_EXT)
BIN_C := $(BIN_DIR)/$(SERVICE_C_NAME)$(BIN_EXT)
BIN_D := $(BIN_DIR)/$(SERVICE_D_NAME)$(BIN_EXT)
BIN_E := $(BIN_DIR)/$(SERVICE_E_NAME)$(BIN_EXT)

# List of all binaries
BINARIES := $(BIN_A) $(BIN_B) $(BIN_C) $(BIN_D) $(BIN_E)

# --- Main Targets ---

.PHONY: all build run clean stop stop-unix

# Target to build and run all
all: build run
 
# 1. Target to build all binaries
build: $(BINARIES)
	@echo "--- All services built successfully in $(BIN_DIR)/ ---"

# 2. Target to run all programs concurrently (Windows Version)
# Uses 'start /B' to run in background and 'pause' to wait
run: build
	@echo "\n--- Running all compiled services concurrently ---"
	@echo "Starting $(SERVICE_A_NAME)..."
	@cmd /c "start /B .\$(subst /,\,$(BIN_A))"
	@echo "Starting $(SERVICE_B_NAME) (Cobra)..."
	@cmd /c "start /B .\$(subst /,\,$(BIN_B)) serve"
	@echo "Starting $(SERVICE_C_NAME)..."
	@cmd /c "start /B .\$(subst /,\,$(BIN_C))"
	@echo "Starting $(SERVICE_D_NAME)..."
	@cmd /c "start /B .\$(subst /,\,$(BIN_D))"
	@echo "Starting $(SERVICE_E_NAME)..."
	@cmd /c "start /B .\$(subst /,\,$(BIN_E))"
	@echo "\nAll services started. Press Enter to stop processes (via 'make stop')."
	@pause

# 4. Target to clean up compiled binaries and directory (Windows Version)
# Uses 'rmdir' instead of 'rm -rf'
clean:
	@echo "--- Cleaning up binaries ---"
	if exist $(subst /,\,$(BIN_DIR)) rmdir /s /q $(subst /,\,$(BIN_DIR))
	@echo "Cleanup complete."

# 5. Target to stop all running services (Windows Version)
# Adds .exe to taskkill and uses '|| exit 0' to ignore errors
stop:
	@echo "--- Stopping all services on Windows (using taskkill) ---"
	@cmd /c "taskkill /F /IM $(SERVICE_A_NAME)$(BIN_EXT) /T 2>nul" || exit 0
	@cmd /c "taskkill /F /IM $(SERVICE_B_NAME)$(BIN_EXT) /T 2>nul" || exit 0
	@cmd /c "taskkill /F /IM $(SERVICE_C_NAME)$(BIN_EXT) /T 2>nul" || exit 0
	@cmd /c "taskkill /F /IM $(SERVICE_D_NAME)$(BIN_EXT) /T 2>nul" || exit 0
	@cmd /c "taskkill /F /IM $(SERVICE_E_NAME)$(BIN_EXT) /T 2>nul" || exit 0
	@echo "All services stopped."

# 6. Target to stop all running services (macOS/Linux Version)
# This target remains unchanged
stop-unix:
	@echo "--- Stopping all services on macOS/Linux (using pkill -f) ---"
	-@pkill -f $(SERVICE_A_NAME) || true
	-@pkill -f $(SERVICE_B_NAME) || true
	-@pkill -f $(SERVICE_C_NAME) || true
	-@pkill -f $(SERVICE_D_NAME) || true
	-@pkill -f $(SERVICE_E_NAME) || true
	@echo "All services stopped."


# --- Utility & Build Rules ---

# 3. Utility to ensure the binary directory exists (Windows Version)
# Uses 'mkdir' instead of 'mkdir -p' and ignores "already exists" error
$(BIN_DIR):
	@echo "--- Creating binary directory: $(BIN_DIR) ---"
	mkdir $(subst /,\,$(BIN_DIR)) 2>nul || true

# Pattern rules for building each binary
# These rules now correctly build the .exe files
$(BIN_A): | $(BIN_DIR)
	@echo "Building $(SERVICE_A_NAME) from $(SERVICE_A_PATH)..."
	go build -o $@ ./$(SERVICE_A_PATH)

$(BIN_B): | $(BIN_DIR)
	@echo "Building $(SERVICE_B_NAME) from $(SERVICE_B_PATH)..."
	go build -o $@ ./$(SERVICE_B_PATH)

$(BIN_C): | $(BIN_DIR)
	@echo "Building $(SERVICE_C_NAME) from $(SERVICE_C_PATH)..."
	go build -o $@ ./$(SERVICE_C_PATH)

$(BIN_D): | $(BIN_DIR)
	@echo "Building $(SERVICE_D_NAME) from $(SERVICE_D_PATH)..."
	go build -o $@ ./$(SERVICE_D_PATH)

$(BIN_E): | $(BIN_DIR)
	@echo "Building $(SERVICE_E_NAME) from $(SERVICE_E_PATH)..."
	go build -o $@ ./$(SERVICE_E_PATH)