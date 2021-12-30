# Build all by default
.DEFAULT_GOAL=all

.PHONY: all
all: build

# ==============================================================================
# Targets

## build: Build source code for host platfor.
.PHONY: build
build:
	@$(MAKE) go.build