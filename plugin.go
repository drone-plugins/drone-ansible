package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var ansibleFolder = "/etc/ansible"
var ansibleConfig = "/etc/ansible/ansible.cfg"

var ansibleContent = `
[defaults]
host_key_checking = False
`

type (
	Config struct {
		Requirements      string
		Galaxy            string
		Inventories       []string
		Playbooks         []string
		Limit             string
		SkipTags          string
		StartAtTask       string
		Tags              string
		ExtraVars         []string
		ModulePath        []string
		Check             bool
		Diff              bool
		FlushCache        bool
		ForceHandlers     bool
		ListHosts         bool
		ListTags          bool
		ListTasks         bool
		SyntaxCheck       bool
		Forks             int
		VaultID           string
		VaultPassword     string
		VaultPasswordFile string
		Verbose           int
		PrivateKey        string
		PrivateKeyFile    string
		User              string
		Connection        string
		Timeout           int
		SSHCommonArgs     string
		SFTPExtraArgs     string
		SCPExtraArgs      string
		SSHExtraArgs      string
		Become            bool
		BecomeMethod      string
		BecomeUser        string
	}

	Plugin struct {
		Config Config
	}
)

func (p *Plugin) Exec() error {
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

func (p *Plugin) ansibleConfig() error {
	if err := os.MkdirAll(ansibleFolder, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create ansible directory")
	}

	if err := ioutil.WriteFile(ansibleConfig, []byte(ansibleContent), 0600); err != nil {
		return errors.Wrap(err, "failed to create ansible config")
	}

	return nil
}

func (p *Plugin) privateKey() error {
	tmpfile, err := ioutil.TempFile("", "privateKey")

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
	tmpfile, err := ioutil.TempFile("", "vaultPass")

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
		"--force",
		"--role-file",
		p.Config.Galaxy,
	}

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

	if p.Config.SyntaxCheck {
		args = append(args, "--syntax-check")
		args = append(args, p.Config.Playbooks...)

		return exec.Command(
			"ansible-playbook",
			args...,
		)
	}

	if p.Config.ListHosts {
		args = append(args, "--list-hosts")
		args = append(args, p.Config.Playbooks...)

		return exec.Command(
			"ansible-playbook",
			args...,
		)
	}

	for _, v := range p.Config.ExtraVars {
		args = append(args, "--extra-vars", v)
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

	if len(p.Config.ModulePath) > 0 {
		args = append(args, "--module-path", strings.Join(p.Config.ModulePath, ":"))
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

	if p.Config.VaultID != "" {
		args = append(args, "--vault-id", p.Config.VaultID)
	}

	if p.Config.VaultPasswordFile != "" {
		args = append(args, "--vault-password-file", p.Config.VaultPasswordFile)
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

	return exec.Command(
		"ansible-playbook",
		args...,
	)
}

func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}
