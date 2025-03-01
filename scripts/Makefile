.PHONY: build install clean uninstall purge check-deps install-deps-arch install-deps-debian

# Colors for terminal output
YELLOW=\033[1;33m
GREEN=\033[1;32m
RED=\033[1;31m
NC=\033[0m # No Color

build:
	@echo "Building LunarFetch..."
	@go build -o lunarfetch

check-deps:
	@echo "Checking dependencies..."
	@missing=""
	@which df free >/dev/null 2>&1 || missing="$$missing coreutils"
	@which uptime >/dev/null 2>&1 || missing="$$missing procps"
	@which lsb_release >/dev/null 2>&1 || missing="$$missing lsb-release"
	@(which xrandr >/dev/null 2>&1 || which xdpyinfo >/dev/null 2>&1 || which swaymsg >/dev/null 2>&1 || which wlr-randr >/dev/null 2>&1) || missing="$$missing display-utils"
	@which lscpu lspci >/dev/null 2>&1 || missing="$$missing pciutils"
	@which gsettings >/dev/null 2>&1 || missing="$$missing gsettings-desktop-schemas"
	@if [ ! -z "$$missing" ]; then \
		echo -e "$(YELLOW)Warning: Missing dependencies:$(NC) $$missing"; \
		echo "Run 'make install-deps-arch' for Arch Linux or 'make install-deps-debian' for Debian/Ubuntu"; \
		echo -e "$(YELLOW)Press Enter to continue anyway, or Ctrl+C to abort...$(NC)"; \
		read dummy; \
	else \
		echo -e "$(GREEN)All dependencies are installed.$(NC)"; \
	fi

install-deps-arch:
	@echo "Installing dependencies for Arch Linux..."
	@sudo pacman -S --needed coreutils procps-ng util-linux lsb-release xorg-xrandr xorg-xdpyinfo pciutils gsettings-desktop-schemas

install-deps-debian:
	@echo "Installing dependencies for Debian/Ubuntu..."
	@sudo apt update && sudo apt install -y coreutils procps lsb-release x11-xserver-utils mesa-utils pciutils gnome-settings-daemon

install: check-deps build
	@echo "Installing LunarFetch..."
	@mkdir -p ~/.config/lunarfetch/logos
	@cp -n assets/config.json ~/.config/lunarfetch/config.json 2>/dev/null || echo "Config file already exists, not overwriting."
	@cp -n assets/logo.txt ~/.config/lunarfetch/logos/moon.txt 2>/dev/null || echo "Logo file already exists, not overwriting."
	@sudo mv lunarfetch /usr/local/bin/
	@echo -e "$(GREEN)LunarFetch has been installed to /usr/local/bin/$(NC)"
	@echo "Run 'lunarfetch' to try it out!"

clean:
	@echo "Cleaning build files..."
	@rm -f lunarfetch
	@echo -e "$(GREEN)Clean complete.$(NC)"

uninstall:
	@echo "Uninstalling LunarFetch binary..."
	@sudo rm -f /usr/local/bin/lunarfetch
	@echo -e "$(GREEN)LunarFetch binary removed.$(NC)"

# Don't remove config files by default
purge: uninstall
	@echo "Removing all LunarFetch files including configuration..."
	@rm -rf ~/.config/lunarfetch
	@echo -e "$(GREEN)All LunarFetch files removed.$(NC)" 