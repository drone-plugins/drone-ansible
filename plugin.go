package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var ansibleFolder = "/etc/ansible"
var ansibleConfig = "/etc/ansible/ansible.cfg"

// var ansibleContent = `
// [defaults]
// host_key_checking = False
// `

const (
	ModePlaybook = "playbook"
	ModeAdhoc    = "adhoc"
	ModeVault    = "vault"
)

// Constants for valid actions
const (
	ActionEncrypt       = "encrypt"
	ActionDecrypt       = "decrypt"
	ActionEncryptString = "encrypt_string"
	ActionView          = "view"
	ActionEdit          = "edit"
	ActionRekey         = "rekey"
)

type (
	Config struct {
		Mode                   string
		Requirements           string
		Galaxy                 string
		Inventories            []string
		Playbooks              []string
		Limit                  string
		SkipTags               string
		StartAtTask            string
		Tags                   string
		ExtraVars              []string
		ModulePath             []string
		GalaxyForce            bool
		Check                  bool
		Diff                   bool
		FlushCache             bool
		ForceHandlers          bool
		ListHosts              bool
		ListTags               bool
		ListTasks              bool
		SyntaxCheck            bool
		Forks                  int
		VaultID                string
		VaultPassword          string
		VaultPasswordFile      string
		Verbose                int
		PrivateKey             string
		PrivateKeyFile         string
		User                   string
		Connection             string
		Timeout                int
		SSHCommonArgs          string
		SFTPExtraArgs          string
		SCPExtraArgs           string
		SSHExtraArgs           string
		Become                 bool
		BecomeMethod           string
		BecomeUser             string
		DisableHostKeyChecking bool   // Disable SSH host key checking
		HostKeyChecking        bool   // Enable SSH host key validation
		Installation           string // Path to the Ansible executable or installation
		InventoryContent       string // Inline inventory content
		Sudo                   bool   // Use sudo for operations
		SudoUser               string // Sudo user for operations
		VaultTmpPath           string // Temporary path for vault password files and others
		// Ad-Hoc Parameters
		Hosts               string // Target hosts for ad-hoc command
		Module              string // Module name for ad-hoc command
		ModuleArguments     string // Module arguments for ad-hoc command
		DynamicInventory    bool   // Enable dynamic inventory
		Extras              string // Additional options for ad-hoc execution
		VaultCredentialsKey string // Vault credentials ID for encrypted files (optional)
		// Inventory          string

		// Vault Parameters
		Action                 string // Action for vault operation (e.g., encrypt, decrypt)
		Content                string // Content for vault operation
		Input                  string // Input file for vault operation
		NewVaultCredentialsKey string // New vault credentials ID for rekeying
		Output                 string // Output file for vault operation
	}

	Plugin struct {
		Config Config
	}
)

func (p *Plugin) Exec() error {
	switch p.Config.Mode {
	case ModePlaybook:
		return p.executePlaybook()
	case ModeAdhoc:
		return p.executeAdhoc()
	case ModeVault:
		return p.executeVault()
	default:
		return errors.New("invalid mode: specify 'playbook' or 'adhoc'")
	}
}

func (p *Plugin) executePlaybook() error {
	if err := p.playbooks(); err != nil {
		return err
	}

	if err := p.ansibleConfig(); err != nil {
		return err
	}

	if p.Config.PrivateKey != "" {
		if err := p.privateKey(); err != nil {
			return err
		}

		defer os.Remove(p.Config.PrivateKeyFile)
	}

	if p.Config.VaultPassword != "" {
		if err := p.vaultPass(); err != nil {
			return err
		}

		defer os.Remove(p.Config.VaultPasswordFile)
	}

	// Handle inline inventory content
	if err := p.setupInventory(); err != nil {
		return err
	}

	// Validate custom Ansible installation
	if err := p.validateInstallation(); err != nil {
		return err
	}

	commands := []*exec.Cmd{
		p.versionCommand(),
	}

	if p.Config.Requirements != "" {
		commands = append(commands, p.requirementsCommand())
	}

	if p.Config.Galaxy != "" {
		commands = append(commands, p.galaxyCommand())
	}

	for _, inventory := range p.Config.Inventories {
		commands = append(commands, p.ansibleCommand(inventory))
	}

	for _, cmd := range commands {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "ANSIBLE_FORCE_COLOR=1")

		trace(cmd)

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

// executeAdhoc executes the Ansible Ad-Hoc command
func (p *Plugin) executeAdhoc() error {
	// Step 1: Validate required parameters
	if p.Config.Hosts == "" {
		return errors.New("hosts parameter is required for ad-hoc execution")
	}

	// Step 2: Default to 'command' module if no module is provided
	module := p.Config.Module
	if module == "" {
		module = "command"
	}

	// Step 3: Build arguments for the ad-hoc command
	args := []string{
		p.Config.Hosts, // Target hosts
		"-m", module,   // Module to execute
	}

	// Step 4: Add module arguments
	if p.Config.ModuleArguments != "" {
		args = append(args, "-a", p.Config.ModuleArguments)
	}

	// Step 5: Add inventory files or inline content
	if len(p.Config.Inventories) > 0 {
		for _, inventory := range p.Config.Inventories {
			args = append(args, "--inventory", inventory)
		}
	}

	if p.Config.InventoryContent != "" {
		tmpfile, err := os.CreateTemp("", "inventory")
		if err != nil {
			return fmt.Errorf("failed to create temporary inventory file: %w", err)
		}
		defer os.Remove(tmpfile.Name())
		if _, err := tmpfile.WriteString(p.Config.InventoryContent); err != nil {
			return fmt.Errorf("failed to write inventory content to temporary file: %w", err)
		}
		args = append(args, "--inventory", tmpfile.Name())
	}

	// Step 6: Handle privilege escalation
	if p.Config.Become {
		args = append(args, "--become")
	}
	if p.Config.BecomeUser != "" {
		args = append(args, "--become-user", p.Config.BecomeUser)
	}

	// Step 7: Handle dynamic inventory
	if p.Config.DynamicInventory {
		args = append(args, "--dynamic-inventory")
	}

	// Step 8: Add extra variables
	for _, ev := range p.Config.ExtraVars {
		args = append(args, "--extra-vars", ev)
	}

	// Step 9: Add additional options
	if p.Config.Extras != "" {
		args = append(args, p.Config.Extras)
	}

	// Step 10: Handle forks for parallelism
	if p.Config.Forks > 0 {
		args = append(args, "--forks", strconv.Itoa(p.Config.Forks))
	}

	// Step 11: Handle host key checking
	env := os.Environ()
	if !p.Config.HostKeyChecking {
		env = append(env, "ANSIBLE_HOST_KEY_CHECKING=False")
	} else {
		env = append(env, "ANSIBLE_HOST_KEY_CHECKING=True")
	}

	// Step 12: Handle vault credentials key
	if p.Config.VaultCredentialsKey != "" {
		tmpVaultFile, err := os.CreateTemp("", "vault-pass")
		if err != nil {
			return fmt.Errorf("failed to create temporary vault password file: %w", err)
		}
		defer os.Remove(tmpVaultFile.Name())
		if _, err := tmpVaultFile.WriteString(p.Config.VaultCredentialsKey); err != nil {
			return fmt.Errorf("failed to write vault password to temporary file: %w", err)
		}
		args = append(args, "--vault-password-file", tmpVaultFile.Name())
	}

	// Step 13: Handle vault temporary path
	if p.Config.VaultTmpPath != "" {
		args = append(args, "--vault-password-file", p.Config.VaultTmpPath)
	}

	// Step 14: Handle private key file
	if p.Config.PrivateKeyFile != "" {
		args = append(args, "--private-key", p.Config.PrivateKeyFile)
	}

	// Step 15: Use custom Ansible installation if provided
	executable := "ansible"
	if p.Config.Installation != "" {
		executable = p.Config.Installation
	}

	// Step 16: Construct and execute the command
	cmd := exec.Command(executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env // Pass environment variables

	// Log the command for debugging purposes
	fmt.Printf("Executing command: %s %v\n", executable, args)

	// Step 17: Run the command
	return cmd.Run()
}

// executeVault executes the Ansible Vault operation
func (p *Plugin) executeVault() error {
	// Step 1: Validate the action
	if err := validateAction(p.Config.Action); err != nil {
		return err
	}

	// Step 2: Determine the ansible-vault executable path
	vaultExecutable := "ansible-vault"
	if p.Config.Installation != "" {
		vaultExecutable = p.Config.Installation
	}

	// Step 3: Build arguments for the ansible-vault command
	args := []string{p.Config.Action}

	// Step 4: Handle input or content based on the action
	if err := handleInputAndContent(p.Config, args); err != nil {
		return err
	}

	// Step 5: Add output file (if applicable)
	handleOutputFile(p.Config.Output, args)

	// Step 6: Handle vault password key
	if err := handleVaultPassword(p.Config.VaultCredentialsKey, args); err != nil {
		return err
	}

	// Step 7: Handle new vault password key for rekeying
	if p.Config.Action == ActionRekey && p.Config.NewVaultCredentialsKey != "" {
		if err := handleNewVaultPassword(p.Config.NewVaultCredentialsKey, args); err != nil {
			return err
		}
	}

	// Step 8: Construct and execute the command
	cmd := exec.Command(vaultExecutable, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Step 9: Log the command for debugging purposes
	fmt.Printf("Executing command: %s %v\n", vaultExecutable, args)

	// Step 10: Execute the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ansible-vault command failed: %w", err)
	}

	return nil
}

// validateAction validates the provided action
func validateAction(action string) error {
	validActions := map[string]bool{
		ActionEncrypt:       true,
		ActionDecrypt:       true,
		ActionEncryptString: true,
		ActionView:          true,
		ActionEdit:          true,
		ActionRekey:         true,
	}
	if !validActions[action] {
		return fmt.Errorf("invalid action: %s. Supported actions: encrypt, decrypt, encrypt_string, view, edit, rekey", action)
	}
	return nil
}

// handleInputAndContent validates and adds input or content based on action
func handleInputAndContent(config Config, args []string) error {
	if config.Action == ActionEncryptString {
		if config.Content == "" {
			return errors.New("content is required for encrypt_string action")
		}

		// If output is provided, do not use --stdin-name
		if config.Output == "" {
			args = append(args, "--stdin-name", "SECRET_VAR")
		}
	} else {
		if config.Input == "" {
			return fmt.Errorf("input file is required for %s action", config.Action)
		}
		args = append(args, config.Input)
	}
	return nil
}

// handleOutputFile adds the output file flag if provided
func handleOutputFile(output string, args []string) {
	if output != "" {
		args = append(args, "--output", output)
	}
}

// handleVaultPassword writes the vault password to a temporary file and appends it to args
func handleVaultPassword(vaultKey string, args []string) error {
	if vaultKey == "" {
		return errors.New("vaultCredentialsKey is required for vault operations")
	}
	tmpVaultFile, err := os.CreateTemp("", "vault-pass")
	if err != nil {
		return fmt.Errorf("failed to create temporary vault password file: %w", err)
	}
	defer os.Remove(tmpVaultFile.Name()) // Ensure cleanup after the function
	if _, err := tmpVaultFile.WriteString(vaultKey); err != nil {
		os.Remove(tmpVaultFile.Name()) // Cleanup on failure
		return fmt.Errorf("failed to write vault key to temporary file: %w", err)
	}
	args = append(args, "--vault-password-file", tmpVaultFile.Name())
	return nil
}

// handleNewVaultPassword writes the new vault password to a temporary file and appends it to args
func handleNewVaultPassword(newVaultKey string, args []string) error {
	tmpNewVaultFile, err := os.CreateTemp("", "new-vault-pass")
	if err != nil {
		return fmt.Errorf("failed to create temporary new vault password file: %w", err)
	}
	defer os.Remove(tmpNewVaultFile.Name()) // Ensure cleanup after the function
	if _, err := tmpNewVaultFile.WriteString(newVaultKey); err != nil {
		os.Remove(tmpNewVaultFile.Name()) // Cleanup on failure
		return fmt.Errorf("failed to write new vault key to temporary file: %w", err)
	}
	args = append(args, "--new-vault-password-file", tmpNewVaultFile.Name())
	return nil
}

// Helper function to locate the vault password file based on the vaultCredentialsId
func locateVaultPasswordFile(vaultCredentialsId string) (string, error) {
	// Simulate retrieval of the vault password file based on the given ID
	// Example: /etc/ansible/vaults/<vaultCredentialsId>.pass
	vaultFilePath := fmt.Sprintf("/etc/ansible/vaults/%s.pass", vaultCredentialsId)
	if _, err := os.Stat(vaultFilePath); os.IsNotExist(err) {
		return "", errors.New("vault password file not found for the given vaultCredentialsId")
	}
	return vaultFilePath, nil
}

func (p *Plugin) ansibleConfig() error {
	if err := os.MkdirAll(ansibleFolder, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create ansible directory")
	}

	ansibleConfigContent := "[defaults]\n"
	if p.Config.DisableHostKeyChecking {
		ansibleConfigContent += "host_key_checking = False\n"
	}

	if err := os.WriteFile(ansibleConfig, []byte(ansibleConfigContent), 0600); err != nil {
		return errors.Wrap(err, "failed to create ansible config")
	}

	return nil
}

func (p *Plugin) privateKey() error {
	tmpfile, err := os.CreateTemp("", "privateKey")

	if err != nil {
		return errors.Wrap(err, "failed to create private key file")
	}

	if _, err := tmpfile.Write([]byte(p.Config.PrivateKey)); err != nil {
		return errors.Wrap(err, "failed to write private key file")
	}

	if err := tmpfile.Close(); err != nil {
		return errors.Wrap(err, "failed to close private key file")
	}

	p.Config.PrivateKeyFile = tmpfile.Name()
	return nil
}

func (p *Plugin) vaultPass() error {
	tmpfile, err := os.CreateTemp("", "vaultPass")

	if err != nil {
		return errors.Wrap(err, "failed to create vault password file")
	}

	if _, err := tmpfile.Write([]byte(p.Config.VaultPassword)); err != nil {
		return errors.Wrap(err, "failed to write vault password file")
	}

	if err := tmpfile.Close(); err != nil {
		return errors.Wrap(err, "failed to close vault password file")
	}

	p.Config.VaultPasswordFile = tmpfile.Name()
	return nil
}

// setupInventory handles inline inventory content
func (p *Plugin) setupInventory() error {
	if p.Config.InventoryContent != "" {
		tmpfile, err := os.CreateTemp("", "inventory")
		if err != nil {
			return errors.Wrap(err, "failed to create temporary inventory file")
		}
		if _, err := tmpfile.WriteString(p.Config.InventoryContent); err != nil {
			return errors.Wrap(err, "failed to write inventory content")
		}
		defer tmpfile.Close()
		p.Config.Inventories = append(p.Config.Inventories, tmpfile.Name())
	}
	return nil
}

// validateInstallation checks if the specified Ansible installation exists
func (p *Plugin) validateInstallation() error {
	if p.Config.Installation != "" {
		if _, err := exec.LookPath(p.Config.Installation); err != nil {
			return errors.Wrapf(err, "specified Ansible installation not found: %s", p.Config.Installation)
		}
	}
	return nil
}

func (p *Plugin) playbooks() error {
	var (
		playbooks []string
	)

	for _, p := range p.Config.Playbooks {
		files, err := filepath.Glob(p)

		if err != nil {
			playbooks = append(playbooks, p)
			continue
		}

		playbooks = append(playbooks, files...)
	}

	if len(playbooks) == 0 {
		return errors.New("failed to find playbook files")
	}

	p.Config.Playbooks = playbooks
	return nil
}

func (p *Plugin) versionCommand() *exec.Cmd {
	args := []string{
		"--version",
	}

	return exec.Command(
		"ansible",
		args...,
	)
}

func (p *Plugin) requirementsCommand() *exec.Cmd {
	args := []string{
		"install",
		"--upgrade",
		"--requirement",
		p.Config.Requirements,
	}

	return exec.Command(
		"pip",
		args...,
	)
}

func (p *Plugin) galaxyCommand() *exec.Cmd {
	args := []string{
		"install",
	}

	if p.Config.GalaxyForce {
		args = append(args, "--force")
	}

	args = append(args,
		"--role-file",
		p.Config.Galaxy,
	)

	if p.Config.Verbose > 0 {
		args = append(args, fmt.Sprintf("-%s", strings.Repeat("v", p.Config.Verbose)))
	}

	return exec.Command(
		"ansible-galaxy",
		args...,
	)
}

func (p *Plugin) ansibleCommand(inventory string) *exec.Cmd {
	args := []string{
		"--inventory",
		inventory,
	}

	if len(p.Config.ModulePath) > 0 {
		args = append(args, "--module-path", strings.Join(p.Config.ModulePath, ":"))
	}

	if p.Config.VaultID != "" {
		args = append(args, "--vault-id", p.Config.VaultID)
	}

	if p.Config.VaultPasswordFile != "" {
		args = append(args, "--vault-password-file", p.Config.VaultPasswordFile)
	}

	if p.Config.VaultTmpPath != "" {
		args = append(args, "--vault-password-file", p.Config.VaultTmpPath) // Vault temporary path
	}

	for _, v := range p.Config.ExtraVars {
		args = append(args, "--extra-vars", v)
	}

	if p.Config.ListHosts {
		args = append(args, "--list-hosts")
		args = append(args, p.Config.Playbooks...)

		return exec.Command(
			"ansible-playbook",
			args...,
		)
	}

	if p.Config.SyntaxCheck {
		args = append(args, "--syntax-check")
		args = append(args, p.Config.Playbooks...)

		return exec.Command(
			"ansible-playbook",
			args...,
		)
	}

	if p.Config.Check {
		args = append(args, "--check")
	}

	if p.Config.Diff {
		args = append(args, "--diff")
	}

	if p.Config.FlushCache {
		args = append(args, "--flush-cache")
	}

	if p.Config.ForceHandlers {
		args = append(args, "--force-handlers")
	}

	if p.Config.Forks != 5 {
		args = append(args, "--forks", strconv.Itoa(p.Config.Forks))
	}

	if p.Config.Limit != "" {
		args = append(args, "--limit", p.Config.Limit)
	}

	if p.Config.ListTags {
		args = append(args, "--list-tags")
	}

	if p.Config.ListTasks {
		args = append(args, "--list-tasks")
	}

	if p.Config.SkipTags != "" {
		args = append(args, "--skip-tags", p.Config.SkipTags)
	}

	if p.Config.StartAtTask != "" {
		args = append(args, "--start-at-task", p.Config.StartAtTask)
	}

	if p.Config.Tags != "" {
		args = append(args, "--tags", p.Config.Tags)
	}

	if p.Config.PrivateKeyFile != "" {
		args = append(args, "--private-key", p.Config.PrivateKeyFile)
	}

	if p.Config.User != "" {
		args = append(args, "--user", p.Config.User)
	}

	if p.Config.Connection != "" {
		args = append(args, "--connection", p.Config.Connection)
	}

	if p.Config.Timeout != 0 {
		args = append(args, "--timeout", strconv.Itoa(p.Config.Timeout))
	}

	if p.Config.SSHCommonArgs != "" {
		args = append(args, "--ssh-common-args", p.Config.SSHCommonArgs)
	}

	if p.Config.SFTPExtraArgs != "" {
		args = append(args, "--sftp-extra-args", p.Config.SFTPExtraArgs)
	}

	if p.Config.SCPExtraArgs != "" {
		args = append(args, "--scp-extra-args", p.Config.SCPExtraArgs)
	}

	if p.Config.SSHExtraArgs != "" {
		args = append(args, "--ssh-extra-args", p.Config.SSHExtraArgs)
	}

	if p.Config.Become {
		args = append(args, "--become")
	}

	if p.Config.BecomeMethod != "" {
		args = append(args, "--become-method", p.Config.BecomeMethod)
	}

	if p.Config.BecomeUser != "" {
		args = append(args, "--become-user", p.Config.BecomeUser)
	}

	if p.Config.Verbose > 0 {
		args = append(args, fmt.Sprintf("-%s", strings.Repeat("v", p.Config.Verbose)))
	}

	args = append(args, p.Config.Playbooks...)

	return exec.Command(p.ansibleExecutable(), args...)

	// return exec.Command(
	// 	"ansible-playbook",
	// 	args...,
	// )
}

// ansibleExecutable determines the executable to use
func (p *Plugin) ansibleExecutable() string {
	if p.Config.Installation != "" {
		return p.Config.Installation
	}
	return "ansible-playbook"
}

func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}
