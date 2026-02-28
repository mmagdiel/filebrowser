package python

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	venvDir         = ".filebrowser_venv"
	requirementsFile = "requirements.txt"
)

// VenvManager manages the Python virtual environment
type VenvManager struct {
	venvPath       string
	pythonPath     string
	pipPath        string
	initialized    bool
}

// NewVenvManager creates a new virtual environment manager
func NewVenvManager() *VenvManager {
	return &VenvManager{
		venvPath: venvDir,
	}
}

// Initialize sets up the virtual environment
func (vm *VenvManager) Initialize() error {
	if vm.initialized {
		return nil
	}

	log.Println("Checking Python virtual environment...")

	// Check if venv already exists
	if vm.venvExists() {
		log.Println("Virtual environment found, verifying...")
		if err := vm.setPaths(); err != nil {
			log.Println("Virtual environment is corrupted, recreating...")
			if err := vm.removeVenv(); err != nil {
				return fmt.Errorf("failed to remove corrupted venv: %w", err)
			}
		} else {
			// Verify dependencies are installed
			if vm.dependenciesInstalled() {
				log.Println("Virtual environment ready")
				vm.initialized = true
				return nil
			}
			log.Println("Dependencies missing, reinstalling...")
		}
	}

	// Create new venv
	log.Println("Creating Python virtual environment...")
	if err := vm.createVenv(); err != nil {
		return fmt.Errorf("failed to create virtual environment: %w", err)
	}

	// Set paths
	if err := vm.setPaths(); err != nil {
		return fmt.Errorf("failed to set venv paths: %w", err)
	}

	// Install dependencies
	log.Println("Installing Python dependencies (this may take a few minutes)...")
	if err := vm.installDependencies(); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	log.Println("Python virtual environment ready")
	vm.initialized = true
	return nil
}

// venvExists checks if the virtual environment directory exists
func (vm *VenvManager) venvExists() bool {
	info, err := os.Stat(vm.venvPath)
	return err == nil && info.IsDir()
}

// createVenv creates a new virtual environment
func (vm *VenvManager) createVenv() error {
	// Try python3 first, then python
	pythonCmd := "python3"
	if _, err := exec.LookPath(pythonCmd); err != nil {
		pythonCmd = "python"
		if _, err := exec.LookPath(pythonCmd); err != nil {
			return fmt.Errorf("python not found in PATH")
		}
	}

	cmd := exec.Command(pythonCmd, "-m", "venv", vm.venvPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("venv creation failed: %v, output: %s", err, string(output))
	}

	return nil
}

// setPaths sets the paths to python and pip in the venv
func (vm *VenvManager) setPaths() error {
	var pythonBin, pipBin string

	if runtime.GOOS == "windows" {
		pythonBin = filepath.Join(vm.venvPath, "Scripts", "python.exe")
		pipBin = filepath.Join(vm.venvPath, "Scripts", "pip.exe")
	} else {
		pythonBin = filepath.Join(vm.venvPath, "bin", "python")
		pipBin = filepath.Join(vm.venvPath, "bin", "pip")
	}

	// Verify python exists
	if _, err := os.Stat(pythonBin); err != nil {
		return fmt.Errorf("python not found in venv: %w", err)
	}

	// Verify pip exists
	if _, err := os.Stat(pipBin); err != nil {
		return fmt.Errorf("pip not found in venv: %w", err)
	}

	vm.pythonPath = pythonBin
	vm.pipPath = pipBin

	return nil
}

// installDependencies installs required Python packages
func (vm *VenvManager) installDependencies() error {
	// Check if requirements.txt exists
	if _, err := os.Stat(requirementsFile); err != nil {
		log.Printf("Warning: %s not found, skipping dependency installation", requirementsFile)
		return nil
	}

	// Upgrade pip first
	log.Println("Upgrading pip...")
	cmd := exec.Command(vm.pipPath, "install", "--upgrade", "pip")
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("Warning: pip upgrade failed: %v, output: %s", err, string(output))
		// Continue anyway, not critical
	}

	// Install requirements
	log.Println("Installing requirements from requirements.txt...")
	cmd = exec.Command(vm.pipPath, "install", "-r", requirementsFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pip install failed: %v, output: %s", err, string(output))
	}

	log.Println("Dependencies installed successfully")
	return nil
}

// dependenciesInstalled checks if required dependencies are installed
func (vm *VenvManager) dependenciesInstalled() bool {
	// Quick check: try to import torch and transformers
	cmd := exec.Command(vm.pythonPath, "-c", "import torch; import transformers")
	err := cmd.Run()
	return err == nil
}

// removeVenv removes the virtual environment directory
func (vm *VenvManager) removeVenv() error {
	return os.RemoveAll(vm.venvPath)
}

// GetPythonPath returns the path to the Python interpreter in the venv
func (vm *VenvManager) GetPythonPath() string {
	return vm.pythonPath
}

// IsInitialized returns whether the venv is initialized
func (vm *VenvManager) IsInitialized() bool {
	return vm.initialized
}

// RunScript runs a Python script using the venv Python interpreter
func (vm *VenvManager) RunScript(scriptPath string, args ...string) error {
	if !vm.initialized {
		return fmt.Errorf("virtual environment not initialized")
	}

	cmdArgs := append([]string{scriptPath}, args...)
	cmd := exec.Command(vm.pythonPath, cmdArgs...)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("script execution failed: %v, output: %s", err, string(output))
	}

	return nil
}
