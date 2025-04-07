package scripts

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	Version     = "1.5.0"
	VersionDate = "2025-03-01"
	BuildDate   = time.Now().Format("2006-01-02")
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

type Dependency struct {
	Name        string
	Commands    []string
	ArchPackage string
	DebPackage  string
}

var dependencies = []Dependency{
	{
		Name:        "coreutils",
		Commands:    []string{"df", "free"},
		ArchPackage: "coreutils",
		DebPackage:  "coreutils",
	},
	{
		Name:        "process utilities",
		Commands:    []string{"uptime"},
		ArchPackage: "procps-ng",
		DebPackage:  "procps",
	},
	{
		Name:        "LSB release",
		Commands:    []string{"lsb_release"},
		ArchPackage: "lsb-release",
		DebPackage:  "lsb-release",
	},
	{
		Name:        "display utilities",
		Commands:    []string{"xrandr", "xdpyinfo", "swaymsg", "wlr-randr"},
		ArchPackage: "xorg-xrandr xorg-xdpyinfo",
		DebPackage:  "x11-xserver-utils",
	},
	{
		Name:        "hardware detection",
		Commands:    []string{"lspci"},
		ArchPackage: "pciutils",
		DebPackage:  "pciutils",
	},
	{
		Name:        "theme detection",
		Commands:    []string{"gsettings", "dconf"},
		ArchPackage: "gsettings-desktop-schemas dconf",
		DebPackage:  "gsettings-desktop-schemas dconf-cli",
	},
	{
		Name:        "image processing",
		Commands:    []string{"chafa", "sixel-convert"},
		ArchPackage: "chafa",
		DebPackage:  "chafa",
	},
}

type SimpleConfig struct {
	Logo Logo `json:"logo"`
	Info Info `json:"info"`
}

type Logo struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
	Path    string `json:"path"`
	Color   string `json:"color"`
}

type Info struct {
	Enabled bool     `json:"enabled"`
	Items   []string `json:"items"`
}

func HandleCommands(args []string) {
	if len(args) == 0 {
		fmt.Printf("%sNo command specified%s\n", ColorRed, ColorReset)
		PrintUsage()
		return
	}

	switch args[0] {
	case "install":
		Install()
	case "uninstall":
		Uninstall(false)
	case "purge":
		Uninstall(true)
	case "check-deps":
		CheckDependencies()
	case "install-deps":
		InstallDependencies()
	case "build":
		BuildBinary()
	case "help", "-h", "--help":
		PrintUsage()
	case "version", "-v", "--version":
		PrintVersion()
	case "setup-image":
		SetupImage()
	default:
		fmt.Printf("%sUnknown command: %s%s\n", ColorRed, args[0], ColorReset)
		PrintUsage()
	}
}

func PrintUsage() {
	fmt.Printf("%sLunarFetch%s - A customizable system information display tool\n\n", ColorCyan, ColorReset)

	printMainUsage()
	printFlagsSection()
	printCommandsSection()
	printExamplesSection()
	printConfigurationSection()

	fmt.Printf("For more information and updates, visit: %shttps://github.com/Lunaris-Project/lunarfetch%s\n", ColorBlue, ColorReset)
}

func printMainUsage() {
	fmt.Printf("%sUSAGE:%s\n", ColorYellow, ColorReset)
	fmt.Printf("  lunarfetch [flags] [command]\n\n")
}

func printFlagsSection() {
	fmt.Printf("%sFLAGS:%s\n", ColorYellow, ColorReset)
	fmt.Printf("  %s-c, --config%s <path>    Specify a custom configuration file path\n", ColorGreen, ColorReset)
	fmt.Printf("  %s-d, --debug%s            Enable debug mode for verbose output\n", ColorGreen, ColorReset)
	fmt.Printf("  %s-h, --help%s             Display this help message\n", ColorGreen, ColorReset)
	fmt.Printf("  %s-v, --version%s          Display version information\n\n", ColorGreen, ColorReset)
}

func printCommandsSection() {
	fmt.Printf("%sCOMMANDS:%s\n", ColorYellow, ColorReset)

	fmt.Printf("  %sinstall%s              Install LunarFetch to your system\n", ColorGreen, ColorReset)
	fmt.Printf("                       - Creates configuration directories\n")
	fmt.Printf("                       - Copies default configuration and assets\n")
	fmt.Printf("                       - Installs binary to /usr/local/bin\n\n")

	fmt.Printf("  %suninstall%s            Remove LunarFetch binary from your system\n", ColorGreen, ColorReset)
	fmt.Printf("                       - Removes only the binary file\n")
	fmt.Printf("                       - Keeps configuration files intact\n\n")

	fmt.Printf("  %spurge%s                Completely remove LunarFetch from your system\n", ColorGreen, ColorReset)
	fmt.Printf("                       - Removes the binary file\n")
	fmt.Printf("                       - Deletes all configuration files and directories\n\n")

	fmt.Printf("  %scheck-deps%s           Check for required system dependencies\n", ColorGreen, ColorReset)
	fmt.Printf("                       - Verifies all required commands are available\n")
	fmt.Printf("                       - Reports missing dependencies\n\n")

	fmt.Printf("  %sinstall-deps%s         Install required system dependencies\n", ColorGreen, ColorReset)
	fmt.Printf("                       - Automatically installs missing dependencies\n")
	fmt.Printf("                       - Supports Arch Linux and Debian/Ubuntu\n\n")

	fmt.Printf("  %sbuild%s                Build the binary without installing\n", ColorGreen, ColorReset)
	fmt.Printf("                       - Creates executable in current directory\n")
	fmt.Printf("                       - Useful for development and testing\n\n")

	fmt.Printf("  %ssetup-image%s          Configure image display support\n", ColorGreen, ColorReset)
	fmt.Printf("                       - Sets up image configuration\n")
	fmt.Printf("                       - Creates necessary directories\n")
	fmt.Printf("                       - Updates configuration file\n\n")

	fmt.Printf("  %shelp%s                 Display this help message\n\n", ColorGreen, ColorReset)

	fmt.Printf("  %sversion%s              Display version information\n\n", ColorGreen, ColorReset)
}

func printExamplesSection() {
	fmt.Printf("%sEXAMPLES:%s\n", ColorYellow, ColorReset)
	fmt.Printf("  lunarfetch                          # Display system information with default config\n")
	fmt.Printf("  lunarfetch -c ~/.config/lunarfetch/custom.json  # Use custom config file\n")
	fmt.Printf("  lunarfetch --debug                  # Run with debug output\n")
	fmt.Printf("  lunarfetch install                  # Install LunarFetch to your system\n")
	fmt.Printf("  lunarfetch setup-image              # Configure image display\n\n")
}

func printConfigurationSection() {
	fmt.Printf("%sCONFIGURATION:%s\n", ColorYellow, ColorReset)
	fmt.Printf("  Default config location: ~/.config/lunarfetch/config.json\n")
	fmt.Printf("  Logo directory: ~/.config/lunarfetch/logos/\n")
	fmt.Printf("  Images directory: ~/.config/lunarfetch/images/\n\n")
}

func PrintVersion() {
	fmt.Printf("%sLunarFetch%s - A modern, customizable system information display tool\n", ColorCyan, ColorReset)
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Release Date: %s\n", VersionDate)
	fmt.Printf("Build Date: %s\n", BuildDate)
	fmt.Printf("For more information, visit: %shttps://github.com/Lunaris-Project/lunarfetch%s\n", ColorBlue, ColorReset)
}

func BuildBinary() bool {
	fmt.Printf("%sBuilding LunarFetch...%s\n", ColorCyan, ColorReset)

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		fmt.Printf("%sInitializing Go module...%s\n", ColorYellow, ColorReset)
		cmd := exec.Command("go", "mod", "init", "lunarfetch")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%sError initializing Go module: %v%s\n", ColorRed, err, ColorReset)
			fmt.Printf("%sOutput: %s%s\n", ColorRed, string(output), ColorReset)
			fmt.Printf("%sTrying to build without modules...%s\n", ColorYellow, ColorReset)
			return false
		} else {

			fmt.Printf("%sTidying Go module...%s\n", ColorYellow, ColorReset)
			cmd = exec.Command("go", "mod", "tidy")
			output, err = cmd.CombinedOutput()
			if err != nil {
				fmt.Printf("%sWarning: Error tidying Go module: %v%s\n", ColorYellow, err, ColorReset)
				fmt.Printf("%sOutput: %s%s\n", ColorYellow, string(output), ColorReset)
				fmt.Printf("%sContinuing with build...%s\n", ColorYellow, ColorReset)
			}
		}
	}

	cmd := exec.Command("go", "build", "-o", "lunarfetch", ".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%sError building binary: %v%s\n", ColorRed, err, ColorReset)
		fmt.Printf("%sOutput: %s%s\n", ColorRed, string(output), ColorReset)

		fmt.Printf("%sTrying alternative build method...%s\n", ColorYellow, ColorReset)
		cmd = exec.Command("go", "build", "-o", "lunarfetch")
		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("%sError building binary: %v%s\n", ColorRed, err, ColorReset)
			fmt.Printf("%sOutput: %s%s\n", ColorRed, string(output), ColorReset)
			return false
		}
	}

	if _, err := os.Stat("lunarfetch"); os.IsNotExist(err) {
		fmt.Printf("%sError: Binary was not created%s\n", ColorRed, ColorReset)
		return false
	}

	fmt.Printf("%sSuccessfully built LunarFetch binary.%s\n", ColorGreen, ColorReset)
	return true
}

func InstallBinary() bool {
	fmt.Printf("%sInstalling LunarFetch binary to /usr/local/bin/...%s\n", ColorCyan, ColorReset)

	absPath, err := filepath.Abs("lunarfetch")
	if err != nil {
		fmt.Printf("%sError getting absolute path: %v%s\n", ColorRed, err, ColorReset)
		return false
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Printf("%sError: Binary not found at %s%s\n", ColorRed, absPath, ColorReset)
		return false
	}

	cmd := exec.Command("sudo", "cp", absPath, "/usr/local/bin/")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("%sError installing binary: %v%s\n", ColorRed, err, ColorReset)

		fmt.Printf("%sTrying alternative installation method...%s\n", ColorYellow, ColorReset)
		cmd = exec.Command("sudo", "install", "-m", "755", absPath, "/usr/local/bin/")
		err = cmd.Run()
		if err != nil {
			fmt.Printf("%sError installing binary: %v%s\n", ColorRed, err, ColorReset)
			return false
		}
	}

	cmd = exec.Command("sudo", "chmod", "+x", "/usr/local/bin/lunarfetch")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("%sWarning: Could not set executable permissions: %v%s\n", ColorYellow, err, ColorReset)

	}

	fmt.Printf("%sLunarFetch has been installed successfully!%s\n", ColorGreen, ColorReset)
	return true
}

func CheckDependencies() []Dependency {
	fmt.Printf("%sChecking for required dependencies...%s\n", ColorCyan, ColorReset)

	var missingDeps []Dependency
	for _, dep := range dependencies {
		if !DependencyExists(dep) {
			missingDeps = append(missingDeps, dep)
			fmt.Printf("  %s✗ %s%s\n", ColorRed, dep.Name, ColorReset)
		} else {
			fmt.Printf("  %s✓ %s%s\n", ColorGreen, dep.Name, ColorReset)
		}
	}

	if len(missingDeps) > 0 {
		fmt.Printf("\n%sSome dependencies are missing.%s\n", ColorYellow, ColorReset)
		fmt.Printf("Run '%slunarfetch install-deps%s' to install them.\n", ColorGreen, ColorReset)
	} else {
		fmt.Printf("\n%sAll dependencies are installed.%s\n", ColorGreen, ColorReset)
	}

	return missingDeps
}

func InstallDependencies() {
	fmt.Printf("%sInstalling dependencies...%s\n", ColorCyan, ColorReset)

	var distro string

	pacmanExists, _ := exec.LookPath("pacman")
	aptExists, _ := exec.LookPath("apt")
	aptGetExists, _ := exec.LookPath("apt-get")

	if pacmanExists != "" {
		distro = "arch"
	} else if aptExists != "" || aptGetExists != "" {
		distro = "debian"
	} else {

		out, err := exec.Command("lsb_release", "-si").Output()
		if err == nil {
			distro = strings.TrimSpace(string(out))
		} else {

			if _, err := os.Stat("/etc/arch-release"); err == nil {
				distro = "arch"
			} else if _, err := os.Stat("/etc/debian_version"); err == nil {
				distro = "debian"
			} else if _, err := os.Stat("/etc/os-release"); err == nil {

				osReleaseContent, err := os.ReadFile("/etc/os-release")
				if err == nil {
					osRelease := string(osReleaseContent)
					if strings.Contains(strings.ToLower(osRelease), "arch") {
						distro = "arch"
					} else if strings.Contains(strings.ToLower(osRelease), "debian") {
						distro = "debian"
					} else if strings.Contains(strings.ToLower(osRelease), "ubuntu") {
						distro = "ubuntu"
					}
				}
			}
		}
	}

	if distro == "" {
		fmt.Printf("%sCould not automatically detect your distribution.%s\n", ColorYellow, ColorReset)
		fmt.Printf("Please select your distribution:\n")
		fmt.Printf("1. Arch Linux (or Arch-based like Manjaro, EndeavourOS, etc.)\n")
		fmt.Printf("2. Debian/Ubuntu (or derivatives)\n")
		fmt.Printf("3. Other\n")
		fmt.Printf("Enter your choice (1-3): ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			distro = "arch"
		case "2":
			distro = "debian"
		default:
			fmt.Printf("%sUnsupported distribution. Please install dependencies manually.%s\n", ColorRed, ColorReset)
			return
		}
	}

	if strings.Contains(strings.ToLower(distro), "arch") || distro == "arch" {
		InstallArchDependencies()
	} else if strings.Contains(strings.ToLower(distro), "debian") || strings.Contains(strings.ToLower(distro), "ubuntu") {
		InstallDebianDependencies()
	} else {
		fmt.Printf("%sUnsupported distribution: %s. Please install dependencies manually.%s\n", ColorRed, distro, ColorReset)
	}
}

func InstallArchDependencies() {
	fmt.Printf("%sInstalling dependencies for Arch Linux...%s\n", ColorCyan, ColorReset)

	var packages []string
	for _, dep := range dependencies {
		if !DependencyExists(dep) && dep.ArchPackage != "" {
			packages = append(packages, strings.Fields(dep.ArchPackage)...)
		}
	}

	if len(packages) == 0 {
		fmt.Printf("%sNo packages to install.%s\n", ColorGreen, ColorReset)
		return
	}

	args := append([]string{"pacman", "-S", "--noconfirm"}, packages...)
	cmd := exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Printf("%sError installing packages: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	fmt.Printf("%sSuccessfully installed dependencies.%s\n", ColorGreen, ColorReset)
}

func InstallDebianDependencies() {
	fmt.Printf("%sInstalling dependencies for Debian/Ubuntu...%s\n", ColorCyan, ColorReset)

	var packages []string
	for _, dep := range dependencies {
		if !DependencyExists(dep) && dep.DebPackage != "" {
			packages = append(packages, strings.Fields(dep.DebPackage)...)
		}
	}

	if len(packages) == 0 {
		fmt.Printf("%sNo packages to install.%s\n", ColorGreen, ColorReset)
		return
	}

	cmd := exec.Command("sudo", "apt", "update")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()

	args := append([]string{"apt", "install", "-y"}, packages...)
	cmd = exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		fmt.Printf("%sError installing packages: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	fmt.Printf("%sSuccessfully installed dependencies.%s\n", ColorGreen, ColorReset)
}

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func DependencyExists(dep Dependency) bool {

	if dep.Name == "LSB release" {

		if _, err := exec.LookPath("pacman"); err == nil {
			for _, cmd := range dep.Commands {
				if cmd == "lsb_release" {

					if _, err := exec.LookPath(cmd); err == nil {
						return true
					}

					return false
				}
			}
		}
	}

	for _, cmd := range dep.Commands {
		if CommandExists(cmd) {
			return true
		}
	}
	return false
}

func FindSourceDirectory() (string, error) {

	if _, err := os.Stat("main.go"); err == nil {

		return ".", nil
	}

	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}", "lunarfetch")
	output, err := cmd.Output()
	if err == nil {
		dir := strings.TrimSpace(string(output))
		if dir != "" {
			return dir, nil
		}
	}

	commonDirs := []string{
		"/usr/local/src/lunarfetch",
		"/usr/src/lunarfetch",
		"/opt/lunarfetch",
	}

	for _, dir := range commonDirs {
		if _, err := os.Stat(filepath.Join(dir, "main.go")); err == nil {
			return dir, nil
		}
	}

	tempDir, err := os.MkdirTemp("", "lunarfetch-install")
	if err == nil {
		fmt.Printf("%sCreating temporary installation directory: %s%s\n", ColorYellow, tempDir, ColorReset)

		modContent := `module lunarfetch

go 1.16
`
		err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(modContent), 0644)
		if err == nil {

			mainContent := `package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	
	cmd := exec.Command("/usr/local/bin/lunarfetch")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}
`
			err = os.WriteFile(filepath.Join(tempDir, "main.go"), []byte(mainContent), 0644)
			if err == nil {
				return tempDir, nil
			}
		}
	}

	return "", fmt.Errorf("could not find LunarFetch source directory")
}

func CopyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Install() {

	sourceDir, err := FindSourceDirectory()
	if err != nil {
		fmt.Printf("%sError: %s%s\n", ColorRed, err.Error(), ColorReset)
		fmt.Printf("%sPlease run this command from the LunarFetch source directory or ensure Go modules are installed.%s\n", ColorYellow, ColorReset)
		return
	}

	fmt.Printf("%sFound LunarFetch source directory: %s%s\n", ColorGreen, sourceDir, ColorReset)

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%sError: Could not get current directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		return
	}

	err = os.Chdir(sourceDir)
	if err != nil {
		fmt.Printf("%sError: Could not change to source directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		return
	}

	defer func() {
		os.Chdir(currentDir)
	}()

	missingDeps := CheckDependencies()
	if len(missingDeps) > 0 {
		fmt.Printf("%sMissing dependencies: %v%s\n", ColorYellow, missingDeps, ColorReset)
		fmt.Printf("%sWould you like to install them? (y/n): %s", ColorYellow, ColorReset)
		var answer string
		fmt.Scanln(&answer)
		if strings.ToLower(answer) == "y" {
			InstallDependencies()
		} else {
			fmt.Printf("%sSkipping dependency installation. Some features may not work correctly.%s\n", ColorYellow, ColorReset)
		}
	}

	fmt.Printf("%sBuilding LunarFetch...%s\n", ColorYellow, ColorReset)
	if BuildBinary() {
		fmt.Printf("%sInstalling LunarFetch...%s\n", ColorYellow, ColorReset)
		InstallBinary()
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("%sError: Could not get home directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		return
	}

	configDir := filepath.Join(homeDir, ".config", "lunarfetch")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		fmt.Printf("%sError: Could not create config directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		return
	}

	logosDir := filepath.Join(configDir, "logos")
	err = os.MkdirAll(logosDir, 0755)
	if err != nil {
		fmt.Printf("%sError: Could not create logos directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		return
	}

	imagesDir := filepath.Join(configDir, "images")
	err = os.MkdirAll(imagesDir, 0755)
	if err != nil {
		fmt.Printf("%sError: Could not create images directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		return
	}

	configPath := filepath.Join(configDir, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		sampleConfigPath := filepath.Join(sourceDir, "src", "assets", "config.json")
		if _, err := os.Stat(sampleConfigPath); err == nil {
			err = CopyFile(sampleConfigPath, configPath)
			if err != nil {
				fmt.Printf("%sError: Could not copy sample config: %s%s\n", ColorRed, err.Error(), ColorReset)
			} else {
				fmt.Printf("%sCreated sample config at %s%s\n", ColorGreen, configPath, ColorReset)
			}
		} else {

			defaultConfig := SimpleConfig{
				Logo: Logo{
					Enabled: true,
					Type:    "ascii",
					Path:    "",
					Color:   "blue",
				},
				Info: Info{
					Enabled: true,
					Items: []string{
						"host", "user", "os", "kernel", "uptime", "terminal", "shell",
						"disk", "memory", "packages", "battery", "gpu", "cpu", "resolution",
						"wm", "theme", "icons", "desktop",
					},
				},
			}
			configJSON, _ := json.MarshalIndent(defaultConfig, "", "  ")
			err = os.WriteFile(configPath, configJSON, 0644)
			if err != nil {
				fmt.Printf("%sError: Could not create default config: %s%s\n", ColorRed, err.Error(), ColorReset)
			} else {
				fmt.Printf("%sCreated default config at %s%s\n", ColorGreen, configPath, ColorReset)
			}
		}
	}

	logoSrcPath := filepath.Join(sourceDir, "src", "assets", "logo.txt")
	if _, err := os.Stat(logoSrcPath); err == nil {
		logoDstPath := filepath.Join(logosDir, "logo.txt")
		err = CopyFile(logoSrcPath, logoDstPath)
		if err != nil {
			fmt.Printf("%sError: Could not copy logo.txt: %s%s\n", ColorRed, err.Error(), ColorReset)
		} else {
			fmt.Printf("%sCreated logo file at %s%s\n", ColorGreen, logoDstPath, ColorReset)
		}
	} else {
		fmt.Printf("%sWarning: Could not find logo.txt in assets directory%s\n", ColorYellow, ColorReset)
	}

	imageSrcPath := filepath.Join(sourceDir, "src", "assets", "image.png")
	if _, err := os.Stat(imageSrcPath); err == nil {
		imageDstPath := filepath.Join(imagesDir, "image.png")
		err = CopyFile(imageSrcPath, imageDstPath)
		if err != nil {
			fmt.Printf("%sError: Could not copy image.png: %s%s\n", ColorRed, err.Error(), ColorReset)
		} else {
			fmt.Printf("%sCreated image file at %s%s\n", ColorGreen, imageDstPath, ColorReset)
		}
	}

	exampleConfigPath := filepath.Join(sourceDir, "example-config.json")
	if _, err := os.Stat(exampleConfigPath); err == nil {
		dstPath := filepath.Join(configDir, "example-config.json")
		err = CopyFile(exampleConfigPath, dstPath)
		if err != nil {
			fmt.Printf("%sError: Could not copy example config: %s%s\n", ColorRed, err.Error(), ColorReset)
		} else {
			fmt.Printf("%sCreated example config at %s%s\n", ColorGreen, dstPath, ColorReset)
		}
	}

	if strings.Contains(sourceDir, "lunarfetch-install") {
		fmt.Printf("%sCleaning up temporary directory...%s\n", ColorYellow, ColorReset)
		os.RemoveAll(sourceDir)
	}

	fmt.Printf("%sLunarFetch has been installed successfully!%s\n", ColorGreen, ColorReset)
	fmt.Printf("%sYou can now run it from anywhere with the command: %slunarfetch%s\n", ColorGreen, ColorCyan, ColorReset)
}

func Uninstall(purge bool) {
	if purge {
		fmt.Printf("%sPurging LunarFetch from your system...%s\n", ColorCyan, ColorReset)
	} else {
		fmt.Printf("%sUninstalling LunarFetch from your system...%s\n", ColorCyan, ColorReset)
	}

	fmt.Printf("Removing LunarFetch binary...\n")
	cmd := exec.Command("sudo", "rm", "-f", "/usr/local/bin/lunarfetch")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%sWarning: Could not remove binary: %s%s\n", ColorYellow, err.Error(), ColorReset)
		fmt.Printf("%sOutput: %s%s\n", ColorYellow, string(output), ColorReset)
	} else {
		fmt.Printf("%sSuccessfully removed LunarFetch binary.%s\n", ColorGreen, ColorReset)
	}

	if purge {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("%sError: Could not find home directory: %s%s\n", ColorRed, err.Error(), ColorReset)
			os.Exit(1)
		}

		configDir := filepath.Join(homeDir, ".config", "lunarfetch")
		fmt.Printf("Removing configuration directory: %s\n", configDir)
		err = os.RemoveAll(configDir)
		if err != nil {
			fmt.Printf("%sWarning: Could not remove configuration directory: %s%s\n", ColorYellow, err.Error(), ColorReset)
		} else {
			fmt.Printf("%sSuccessfully removed configuration directory.%s\n", ColorGreen, ColorReset)
		}

		tempDir := filepath.Join(os.TempDir(), "lunarfetch-*")
		fmt.Printf("Removing temporary files: %s\n", tempDir)
		cmd = exec.Command("rm", "-rf", tempDir)
		cmd.Run()
	}

	if purge {
		fmt.Printf("%sLunarFetch has been completely purged from your system.%s\n", ColorGreen, ColorReset)
	} else {
		fmt.Printf("%sLunarFetch has been uninstalled from your system.%s\n", ColorGreen, ColorReset)
		fmt.Printf("Your configuration files have been kept at ~/.config/lunarfetch/\n")
		fmt.Printf("To remove them as well, use: %slunarfetch purge%s\n", ColorCyan, ColorReset)
	}
}

func SetupImage() {
	fmt.Printf("%sConfiguring image display support...%s\n", ColorCyan, ColorReset)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("%sError: Could not get home directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		os.Exit(1)
	}

	configDir := filepath.Join(homeDir, ".config", "lunarfetch")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		fmt.Printf("%sError: Could not create config directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		os.Exit(1)
	}

	imagesDir := filepath.Join(configDir, "images")
	err = os.MkdirAll(imagesDir, 0755)
	if err != nil {
		fmt.Printf("%sError: Could not create images directory: %s%s\n", ColorRed, err.Error(), ColorReset)
		os.Exit(1)
	}

	configFile := filepath.Join(configDir, "config.json")

	type Config struct {
		Decorations struct {
			TopLeft     string `json:"topLeft"`
			TopRight    string `json:"topRight"`
			BottomLeft  string `json:"bottomLeft"`
			BottomRight string `json:"bottomRight"`
			TopEdge     string `json:"topEdge"`
			BottomEdge  string `json:"bottomEdge"`
			LeftEdge    string `json:"leftEdge"`
			RightEdge   string `json:"rightEdge"`
			Separator   string `json:"separator"`
		} `json:"decorations"`
		Logo struct {
			EnableLogo bool   `json:"enableLogo"`
			Type       string `json:"type"`
			Content    string `json:"content"`
			Location   string `json:"location"`
			LogoPath   string `json:"logoPath"`
			Position   string `json:"position"`
		} `json:"logo"`
		Image struct {
			EnableImage    bool   `json:"enableImage"`
			Enabled        bool   `json:"enabled"`
			Random         bool   `json:"random"`
			ImagePath      string `json:"imagePath"`
			Width          int    `json:"width"`
			Height         int    `json:"height"`
			RenderMode     string `json:"renderMode"`
			DitherMode     string `json:"ditherMode"`
			TerminalOutput bool   `json:"terminalOutput"`
			DisplayMode    string `json:"displayMode"`
			Protocol       string `json:"protocol"`
			Scale          int    `json:"scale"`
			Offset         int    `json:"offset"`
			Background     string `json:"background"`
			Position       string `json:"position"`
		} `json:"image"`
		Display struct {
			ShowLogoFirst  bool `json:"showLogoFirst"`
			ShowImageFirst bool `json:"showImageFirst"`
		} `json:"display"`
		Icons   map[string]string `json:"icons"`
		Modules map[string]bool   `json:"modules"`
	}

	var config Config

	if _, err := os.Stat(configFile); os.IsNotExist(err) {

		config.Logo.EnableLogo = true
		config.Logo.Type = "file"
		config.Logo.LogoPath = filepath.Join(configDir, "logos")

		config.Image.EnableImage = true
		config.Image.Enabled = true
		config.Image.ImagePath = imagesDir
		config.Image.Random = true
		config.Image.Width = 80
		config.Image.Height = 24
		config.Image.RenderMode = "detailed"
		config.Image.DitherMode = "none"
	} else {

		configData, err := os.ReadFile(configFile)
		if err != nil {
			fmt.Printf("%sError: Could not read config file: %s%s\n", ColorRed, err.Error(), ColorReset)
			os.Exit(1)
		}

		err = json.Unmarshal(configData, &config)
		if err != nil {
			fmt.Printf("%sError: Could not parse config file: %s%s\n", ColorRed, err.Error(), ColorReset)
			fmt.Printf("Creating a new config file...\n")

			config.Logo.EnableLogo = true
			config.Logo.Type = "file"
			config.Logo.LogoPath = filepath.Join(configDir, "logos")

			config.Image.EnableImage = true
			config.Image.Enabled = true
			config.Image.ImagePath = imagesDir
			config.Image.Random = true
			config.Image.Width = 80
			config.Image.Height = 24
			config.Image.RenderMode = "detailed"
			config.Image.DitherMode = "none"
		}
	}

	fmt.Printf("\n%sImage Configuration Options:%s\n", ColorYellow, ColorReset)
	fmt.Printf("1. Add a new image\n")
	fmt.Printf("2. Enable random image selection\n")
	fmt.Printf("3. Set specific image\n")
	fmt.Printf("4. List available images\n")
	fmt.Printf("5. Exit\n")
	fmt.Printf("\nEnter your choice (1-5): ")

	var choice string
	fmt.Scanln(&choice)

	switch choice {
	case "1":

		fmt.Printf("\nPlease enter the path to an image file: ")
		var imagePath string
		fmt.Scanln(&imagePath)

		if imagePath != "" {

			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				fmt.Printf("%sError: Image file not found: %s%s\n", ColorRed, imagePath, ColorReset)
				os.Exit(1)
			}

			fileName := filepath.Base(imagePath)
			destPath := filepath.Join(imagesDir, fileName)
			err = CopyFile(imagePath, destPath)
			if err != nil {
				fmt.Printf("%sError: Could not copy image file: %s%s\n", ColorRed, err.Error(), ColorReset)
				os.Exit(1)
			}

			fmt.Printf("%sImage file copied to: %s%s\n", ColorGreen, destPath, ColorReset)

			fmt.Printf("\nDo you want to use this image specifically or enable random selection?\n")
			fmt.Printf("1. Use this image specifically\n")
			fmt.Printf("2. Enable random selection\n")
			fmt.Printf("Enter your choice (1-2): ")

			var imgChoice string
			fmt.Scanln(&imgChoice)

			if imgChoice == "1" {
				config.Image.EnableImage = true
				config.Image.Enabled = true
				config.Image.Random = false
				config.Image.ImagePath = destPath
			} else {
				config.Image.EnableImage = true
				config.Image.Enabled = true
				config.Image.Random = true
				config.Image.ImagePath = imagesDir
			}
		}
	case "2":

		config.Image.EnableImage = true
		config.Image.Enabled = true
		config.Image.Random = true
		config.Image.ImagePath = imagesDir
		fmt.Printf("%sRandom image selection enabled.%s\n", ColorGreen, ColorReset)

		files, err := os.ReadDir(imagesDir)
		if err != nil {
			fmt.Printf("%sWarning: Could not read images directory: %s%s\n", ColorYellow, err.Error(), ColorReset)
		} else {
			var imageCount int
			for _, file := range files {
				if !file.IsDir() && (strings.HasSuffix(strings.ToLower(file.Name()), ".png") ||
					strings.HasSuffix(strings.ToLower(file.Name()), ".jpg") ||
					strings.HasSuffix(strings.ToLower(file.Name()), ".jpeg")) {
					imageCount++
				}
			}

			if imageCount == 0 {
				fmt.Printf("%sWarning: No image files found in %s. Please add some images.%s\n", ColorYellow, imagesDir, ColorReset)
			} else {
				fmt.Printf("%sFound %d image(s) in the directory.%s\n", ColorGreen, imageCount, ColorReset)
			}
		}
	case "3":

		files, err := os.ReadDir(imagesDir)
		if err != nil {
			fmt.Printf("%sError: Could not read images directory: %s%s\n", ColorRed, err.Error(), ColorReset)
			os.Exit(1)
		}

		var images []string
		fmt.Printf("\n%sAvailable images:%s\n", ColorYellow, ColorReset)
		for i, file := range files {
			if !file.IsDir() && (strings.HasSuffix(strings.ToLower(file.Name()), ".png") ||
				strings.HasSuffix(strings.ToLower(file.Name()), ".jpg") ||
				strings.HasSuffix(strings.ToLower(file.Name()), ".jpeg")) {
				images = append(images, file.Name())
				fmt.Printf("%d. %s\n", i+1, file.Name())
			}
		}

		if len(images) == 0 {
			fmt.Printf("%sNo image files found. Please add some images first.%s\n", ColorRed, ColorReset)
			os.Exit(1)
		}

		fmt.Printf("\nEnter the number of the image to use: ")
		var imgIndex int
		fmt.Scanln(&imgIndex)

		if imgIndex < 1 || imgIndex > len(images) {
			fmt.Printf("%sInvalid selection.%s\n", ColorRed, ColorReset)
			os.Exit(1)
		}

		selectedImage := filepath.Join(imagesDir, images[imgIndex-1])
		config.Image.EnableImage = true
		config.Image.Enabled = true
		config.Image.Random = false
		config.Image.ImagePath = selectedImage

		fmt.Printf("%sSelected image: %s%s\n", ColorGreen, selectedImage, ColorReset)
	case "4":

		files, err := os.ReadDir(imagesDir)
		if err != nil {
			fmt.Printf("%sError: Could not read images directory: %s%s\n", ColorRed, err.Error(), ColorReset)
			os.Exit(1)
		}

		fmt.Printf("\n%sAvailable images:%s\n", ColorYellow, ColorReset)
		var imageCount int
		for _, file := range files {
			if !file.IsDir() && (strings.HasSuffix(strings.ToLower(file.Name()), ".png") ||
				strings.HasSuffix(strings.ToLower(file.Name()), ".jpg") ||
				strings.HasSuffix(strings.ToLower(file.Name()), ".jpeg")) {
				fmt.Printf("- %s\n", file.Name())
				imageCount++
			}
		}

		if imageCount == 0 {
			fmt.Printf("%sNo image files found.%s\n", ColorRed, ColorReset)
		}

		SetupImage()
		return
	case "5":

		fmt.Printf("%sExiting image configuration.%s\n", ColorYellow, ColorReset)
		return
	default:
		fmt.Printf("%sInvalid choice. Exiting.%s\n", ColorRed, ColorReset)
		return
	}

	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Printf("%sError: Could not generate config data: %s%s\n", ColorRed, err.Error(), ColorReset)
		os.Exit(1)
	}

	err = os.WriteFile(configFile, configData, 0644)
	if err != nil {
		fmt.Printf("%sError: Could not write config file: %s%s\n", ColorRed, err.Error(), ColorReset)
		os.Exit(1)
	}

	fmt.Printf("%sImage display support has been configured successfully!%s\n", ColorGreen, ColorReset)
	fmt.Printf("Configuration saved to: %s\n", configFile)
	fmt.Printf("\nYou can now run %slunarfetch%s to see your system information with the configured image.\n", ColorCyan, ColorReset)
}
