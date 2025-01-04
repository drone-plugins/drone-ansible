package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var (
	version = "unknown"
)

func main() {
	app := cli.NewApp()
	app.Name = "ansible plugin"
	app.Usage = "ansible plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "mode",
			Usage:  "Mode of the functionality",
			EnvVar: "PLUGIN_MODE",
		},
		cli.StringFlag{
			Name:   "requirements",
			Usage:  "path to python requirements",
			EnvVar: "PLUGIN_REQUIREMENTS",
		},
		cli.StringFlag{
			Name:   "galaxy",
			Usage:  "path to galaxy requirements",
			EnvVar: "PLUGIN_GALAXY",
		},
		cli.StringSliceFlag{
			Name:   "inventory",
			Usage:  "specify inventory host path",
			EnvVar: "PLUGIN_INVENTORY,PLUGIN_INVENTORIES",
		},
		cli.StringSliceFlag{
			Name:   "playbook",
			Usage:  "list of playbooks to apply",
			EnvVar: "PLUGIN_PLAYBOOK,PLUGIN_PLAYBOOKS",
		},
		cli.StringFlag{
			Name:   "limit",
			Usage:  "further limit selected hosts to an additional pattern",
			EnvVar: "PLUGIN_LIMIT",
		},
		cli.StringFlag{
			Name:   "skip-tags",
			Usage:  "only run plays and tasks whose tags do not match",
			EnvVar: "PLUGIN_SKIP_TAGS",
		},
		cli.StringFlag{
			Name:   "start-at-task",
			Usage:  "start the playbook at the task matching this name",
			EnvVar: "PLUGIN_START_AT_TASK",
		},
		cli.StringFlag{
			Name:   "tags",
			Usage:  "only run plays and tasks tagged with these values",
			EnvVar: "PLUGIN_TAGS",
		},
		cli.StringSliceFlag{
			Name:   "extra-vars",
			Usage:  "set additional variables as key=value",
			EnvVar: "PLUGIN_EXTRA_VARS,ANSIBLE_EXTRA_VARS",
		},
		cli.StringSliceFlag{
			Name:   "module-path",
			Usage:  "prepend paths to module library",
			EnvVar: "PLUGIN_MODULE_PATH",
		},
		cli.BoolTFlag{
			Name:   "galaxy-force",
			Usage:  "force overwriting an existing role or collection",
			EnvVar: "PLUGIN_GALAXY_FORCE",
		},
		cli.BoolFlag{
			Name:   "check",
			Usage:  "run a check, do not apply any changes",
			EnvVar: "PLUGIN_CHECK",
		},
		cli.BoolFlag{
			Name:   "diff",
			Usage:  "show the differences, may print secrets",
			EnvVar: "PLUGIN_DIFF",
		},
		cli.BoolFlag{
			Name:   "flush-cache",
			Usage:  "clear the fact cache for every host in inventory",
			EnvVar: "PLUGIN_FLUSH_CACHE",
		},
		cli.BoolFlag{
			Name:   "force-handlers",
			Usage:  "run handlers even if a task fails",
			EnvVar: "PLUGIN_FORCE_HANDLERS",
		},
		cli.BoolFlag{
			Name:   "list-hosts",
			Usage:  "outputs a list of matching hosts",
			EnvVar: "PLUGIN_LIST_HOSTS",
		},
		cli.BoolFlag{
			Name:   "list-tags",
			Usage:  "list all available tags",
			EnvVar: "PLUGIN_LIST_TAGS",
		},
		cli.BoolFlag{
			Name:   "list-tasks",
			Usage:  "list all tasks that would be executed",
			EnvVar: "PLUGIN_LIST_TASKS",
		},
		cli.BoolFlag{
			Name:   "syntax-check",
			Usage:  "perform a syntax check on the playbook",
			EnvVar: "PLUGIN_SYNTAX_CHECK",
		},
		cli.IntFlag{
			Name:   "forks",
			Usage:  "specify number of parallel processes to use",
			EnvVar: "PLUGIN_FORKS",
			Value:  5,
		},
		cli.StringFlag{
			Name:   "vault-id",
			Usage:  "the vault identity to use",
			EnvVar: "PLUGIN_VAULT_ID,ANSIBLE_VAULT_ID",
		},
		cli.StringFlag{
			Name:   "vault-password",
			Usage:  "the vault password to use",
			EnvVar: "PLUGIN_VAULT_PASSWORD,ANSIBLE_VAULT_PASSWORD",
		},
		cli.IntFlag{
			Name:   "verbose",
			Usage:  "level of verbosity, 0 up to 4",
			EnvVar: "PLUGIN_VERBOSE",
		},
		cli.StringFlag{
			Name:   "private-key",
			Usage:  "use this key to authenticate the connection",
			EnvVar: "PLUGIN_PRIVATE_KEY,ANSIBLE_PRIVATE_KEY",
		},
		cli.StringFlag{
			Name:   "user",
			Usage:  "connect as this user",
			EnvVar: "PLUGIN_USER,ANSIBLE_USER",
		},
		cli.StringFlag{
			Name:   "connection",
			Usage:  "connection type to use",
			EnvVar: "PLUGIN_CONNECTION",
		},
		cli.IntFlag{
			Name:   "timeout",
			Usage:  "override the connection timeout in seconds",
			EnvVar: "PLUGIN_TIMEOUT",
		},
		cli.StringFlag{
			Name:   "ssh-common-args",
			Usage:  "specify common arguments to pass to sftp/scp/ssh",
			EnvVar: "PLUGIN_SSH_COMMON_ARGS",
		},
		cli.StringFlag{
			Name:   "sftp-extra-args",
			Usage:  "specify extra arguments to pass to sftp only",
			EnvVar: "PLUGIN_SFTP_EXTRA_ARGS",
		},
		cli.StringFlag{
			Name:   "scp-extra-args",
			Usage:  "specify extra arguments to pass to scp only",
			EnvVar: "PLUGIN_SCP_EXTRA_ARGS",
		},
		cli.StringFlag{
			Name:   "ssh-extra-args",
			Usage:  "specify extra arguments to pass to ssh only",
			EnvVar: "PLUGIN_SSH_EXTRA_ARGS",
		},
		cli.BoolFlag{
			Name:   "become",
			Usage:  "run operations with become",
			EnvVar: "PLUGIN_BECOME",
		},
		cli.StringFlag{
			Name:   "become-method",
			Usage:  "privilege escalation method to use",
			EnvVar: "PLUGIN_BECOME_METHOD,ANSIBLE_BECOME_METHOD",
		},
		cli.StringFlag{
			Name:   "become-user",
			Usage:  "run operations as this user",
			EnvVar: "PLUGIN_BECOME_USER,ANSIBLE_BECOME_USER",
		},
		cli.BoolFlag{
			Name:   "disable-host-key-checking",
			Usage:  "Disable validation of the host's SSH server keys",
			EnvVar: "PLUGIN_DISABLE_HOST_KEY_CHECKING",
		},
		cli.BoolFlag{
			Name:   "host-key-checking",
			Usage:  "Enable validation of the host's SSH server keys",
			EnvVar: "PLUGIN_HOST_KEY_CHECKING",
		},
		cli.StringFlag{
			Name:   "installation",
			Usage:  "Specify the path to Ansible installation",
			EnvVar: "PLUGIN_INSTALLATION",
		},
		cli.StringFlag{
			Name:   "inventory-content",
			Usage:  "Inline inventory content as a string",
			EnvVar: "PLUGIN_INVENTORY_CONTENT",
		},
		cli.BoolFlag{
			Name:   "sudo",
			Usage:  "Use sudo for operations",
			EnvVar: "PLUGIN_SUDO",
		},
		cli.StringFlag{
			Name:   "sudo-user",
			Usage:  "Specify the sudo user (default: root)",
			EnvVar: "PLUGIN_SUDO_USER",
		},
		cli.StringFlag{
			Name:   "vault-tmp-path",
			Usage:  "Temporary path for generated vault files",
			EnvVar: "PLUGIN_VAULT_TMP_PATH",
		},
		// Ad-Hoc Specific Flags
		cli.StringFlag{
			Name:   "hosts",
			Usage:  "Target hosts for ad-hoc command",
			EnvVar: "PLUGIN_HOSTS",
		},
		cli.StringFlag{
			Name:   "module",
			Usage:  "Module name for ad-hoc execution",
			EnvVar: "PLUGIN_MODULE",
		},
		cli.StringFlag{
			Name:   "module-arguments",
			Usage:  "Arguments for the specified module",
			EnvVar: "PLUGIN_MODULE_ARGUMENTS",
		},
		cli.BoolFlag{
			Name:   "dynamic-inventory",
			Usage:  "Enable dynamic inventory",
			EnvVar: "PLUGIN_DYNAMIC_INVENTORY",
		},
		cli.StringFlag{
			Name:   "extras",
			Usage:  "Additional options for ad-hoc execution",
			EnvVar: "PLUGIN_EXTRAS",
		},
		cli.StringFlag{
			Name:   "vault-credentials-key",
			Usage:  "Vault credentials ID for encrypted files",
			EnvVar: "PLUGIN_VAULT_CREDENTIALS_KEY",
		},
		// Vault Specific Flags
		cli.StringFlag{
			Name:   "action",
			Usage:  "Action for ansible-vault (e.g., encrypt, decrypt, view)",
			EnvVar: "PLUGIN_ACTION",
		},
		cli.StringFlag{
			Name:   "content",
			Usage:  "Content to encrypt or decrypt",
			EnvVar: "PLUGIN_CONTENT",
		},
		cli.StringFlag{
			Name:   "input",
			Usage:  "Input file for the vault operation",
			EnvVar: "PLUGIN_INPUT",
		},
		cli.StringFlag{
			Name:   "output",
			Usage:  "Output file for the vault operation",
			EnvVar: "PLUGIN_OUTPUT",
		},
		cli.StringFlag{
			Name:   "new-vault-credentials-key",
			Usage:  "New Vault Credentials Key for rekeying",
			EnvVar: "PLUGIN_NEW_VAULT_CREDENTIALS_KEY",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Config: Config{
			Mode:                   c.String("mode"),
			Requirements:           c.String("requirements"),
			Galaxy:                 c.String("galaxy"),
			Inventories:            c.StringSlice("inventory"),
			Playbooks:              c.StringSlice("playbook"),
			Limit:                  c.String("limit"),
			SkipTags:               c.String("skip-tags"),
			StartAtTask:            c.String("start-at-task"),
			Tags:                   c.String("tags"),
			ExtraVars:              c.StringSlice("extra-vars"),
			ModulePath:             c.StringSlice("module-path"),
			GalaxyForce:            c.Bool("galaxy-force"),
			Check:                  c.Bool("check"),
			Diff:                   c.Bool("diff"),
			FlushCache:             c.Bool("flush-cache"),
			ForceHandlers:          c.Bool("force-handlers"),
			ListHosts:              c.Bool("list-hosts"),
			ListTags:               c.Bool("list-tags"),
			ListTasks:              c.Bool("list-tasks"),
			SyntaxCheck:            c.Bool("syntax-check"),
			Forks:                  c.Int("forks"),
			VaultID:                c.String("vailt-id"),
			VaultPassword:          c.String("vault-password"),
			Verbose:                c.Int("verbose"),
			PrivateKey:             c.String("private-key"),
			User:                   c.String("user"),
			Connection:             c.String("connection"),
			Timeout:                c.Int("timeout"),
			SSHCommonArgs:          c.String("ssh-common-args"),
			SFTPExtraArgs:          c.String("sftp-extra-args"),
			SCPExtraArgs:           c.String("scp-extra-args"),
			SSHExtraArgs:           c.String("ssh-extra-args"),
			Become:                 c.Bool("become"),
			BecomeMethod:           c.String("become-method"),
			BecomeUser:             c.String("become-user"),
			DisableHostKeyChecking: c.Bool("disable-host-key-checking"), // Disable SSH host key checking
			HostKeyChecking:        c.Bool("host-key-checking"),         // Enable SSH host key validation
			Installation:           c.String("installation"),            // Path to the Ansible executable or installation
			InventoryContent:       c.String("inventory-content"),       // Inline inventory content
			Sudo:                   c.Bool("sudo"),                      // Use sudo for operations
			SudoUser:               c.String("sudo-user"),               // Sudo user for operations
			VaultTmpPath:           c.String("vault-tmp-path"),          // Temporary path for vault password files and others
			// Ad-Hoc Parameters
			Hosts:               c.String("hosts"),                 // Target hosts for ad-hoc command
			Module:              c.String("module"),                // Module name for ad-hoc command
			ModuleArguments:     c.String("module-arguments"),      // Module arguments for ad-hoc command
			DynamicInventory:    c.Bool("dynamic-inventory"),       // Enable dynamic inventory
			Extras:              c.String("extras"),                // Additional options for ad-hoc execution
			VaultCredentialsKey: c.String("vault-credentials-key"), // Vault credentials ID for encrypted files
			// Vault Parameters
			Action:                 c.String("action"),
			Content:                c.String("content"),
			Input:                  c.String("input"),
			Output:                 c.String("output"),
			NewVaultCredentialsKey: c.String("new-vault-credentials-key"),
		},
	}

	// Set default mode to "playbook" if not explicitly provided
	if plugin.Config.Mode == "" {
		plugin.Config.Mode = "playbook"
	}

	// Validate mode and required parameters based on the mode
	switch plugin.Config.Mode {
	case "playbook":
		if len(plugin.Config.Playbooks) == 0 {
			return errors.New("you must provide a playbook in playbook mode")
		}
		if len(plugin.Config.Inventories) == 0 && plugin.Config.InventoryContent == "" {
			return errors.New("you must provide an inventory or inventory content in playbook mode")
		}
	case "adhoc":
		if plugin.Config.Hosts == "" {
			return errors.New("you must provide hosts for adhoc mode")
		}
		// Module is optional; defaults to "command" if not provided
	case "vault":
		if plugin.Config.VaultCredentialsKey == "" {
			return errors.New("VaultCredentialsKey is mandatory for vault mode")
		}
		// Action, Content, Input, and Output are optional
	default:
		return errors.New("invalid mode: specify 'playbook', 'adhoc', or 'vault'")
	}

	return plugin.Exec()
}
