package main

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var (
	version = "unknown"
)

func main() {

	app := &cli.App{
		Name:    "ansible plugin",
		Usage:   "ansible plugin",
		Version: version,
		Action:  run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "requirements",
				Usage:   "path to python requirements",
				EnvVars: []string{"PLUGIN_REQUIREMENTS"},
			},
			&cli.StringFlag{
				Name:    "galaxy",
				Usage:   "path to galaxy requirements",
				EnvVars: []string{"PLUGIN_GALAXY"},
			},
			&cli.StringSliceFlag{
				Name:    "inventory",
				Usage:   "specify inventory host path",
				EnvVars: []string{"PLUGIN_INVENTORY", "PLUGIN_INVENTORIES"},
			},
			&cli.StringSliceFlag{
				Name:    "playbook",
				Usage:   "list of playbooks to apply",
				EnvVars: []string{"PLUGIN_PLAYBOOK", "PLUGIN_PLAYBOOKS"},
			},
			&cli.StringFlag{
				Name:    "limit",
				Usage:   "further limit selected hosts to an additional pattern",
				EnvVars: []string{"PLUGIN_LIMIT"},
			},
			&cli.StringFlag{
				Name:    "skip-tags",
				Usage:   "only run plays and tasks whose tags do not match",
				EnvVars: []string{"PLUGIN_SKIP_TAGS"},
			},
			&cli.StringFlag{
				Name:    "start-at-task",
				Usage:   "start the playbook at the task matching this name",
				EnvVars: []string{"PLUGIN_START_AT_TASK"},
			},
			&cli.StringFlag{
				Name:    "tags",
				Usage:   "only run plays and tasks tagged with these values",
				EnvVars: []string{"PLUGIN_TAGS"},
			},
			&cli.StringSliceFlag{
				Name:    "extra-vars",
				Usage:   "set additional variables as key=value",
				EnvVars: []string{"PLUGIN_EXTRA_VARS", "ANSIBLE_EXTRA_VARS"},
			},
			&cli.StringSliceFlag{
				Name:    "module-path",
				Usage:   "prepend paths to module library",
				EnvVars: []string{"PLUGIN_MODULE_PATH"},
			},
			&cli.BoolFlag{
				Name:    "galaxy-force",
				Usage:   "force overwriting an existing role or collection",
				EnvVars: []string{"PLUGIN_GALAXY_FORCE"},
			},
			&cli.BoolFlag{
				Name:    "check",
				Usage:   "run a check, do not apply any changes",
				EnvVars: []string{"PLUGIN_CHECK"},
			},
			&cli.BoolFlag{
				Name:    "diff",
				Usage:   "show the differences, may print secrets",
				EnvVars: []string{"PLUGIN_DIFF"},
			},
			&cli.BoolFlag{
				Name:    "flush-cache",
				Usage:   "clear the fact cache for every host in inventory",
				EnvVars: []string{"PLUGIN_FLUSH_CACHE"},
			},
			&cli.BoolFlag{
				Name:    "force-handlers",
				Usage:   "run handlers even if a task fails",
				EnvVars: []string{"PLUGIN_FORCE_HANDLERS"},
			},
			&cli.BoolFlag{
				Name:    "list-hosts",
				Usage:   "outputs a list of matching hosts",
				EnvVars: []string{"PLUGIN_LIST_HOSTS"},
			},
			&cli.BoolFlag{
				Name:    "list-tags",
				Usage:   "list all available tags",
				EnvVars: []string{"PLUGIN_LIST_TAGS"},
			},
			&cli.BoolFlag{
				Name:    "list-tasks",
				Usage:   "list all tasks that would be executed",
				EnvVars: []string{"PLUGIN_LIST_TASKS"},
			},
			&cli.BoolFlag{
				Name:    "syntax-check",
				Usage:   "perform a syntax check on the playbook",
				EnvVars: []string{"PLUGIN_SYNTAX_CHECK"},
			},
			&cli.IntFlag{
				Name:    "forks",
				Usage:   "specify number of parallel processes to use",
				EnvVars: []string{"PLUGIN_FORKS"},
				Value:   5,
			},
			&cli.StringFlag{
				Name:    "vault-id",
				Usage:   "the vault identity to use",
				EnvVars: []string{"PLUGIN_VAULT_ID", "ANSIBLE_VAULT_ID"},
			},
			&cli.StringFlag{
				Name:    "vault-password",
				Usage:   "the vault password to use",
				EnvVars: []string{"PLUGIN_VAULT_PASSWORD", "ANSIBLE_VAULT_PASSWORD"},
			},
			&cli.IntFlag{
				Name:    "verbose",
				Usage:   "level of verbosity, 0 up to 4",
				EnvVars: []string{"PLUGIN_VERBOSE"},
			},
			&cli.StringFlag{
				Name:    "private-key",
				Usage:   "use this key to authenticate the connection",
				EnvVars: []string{"PLUGIN_PRIVATE_KEY", "ANSIBLE_PRIVATE_KEY"},
			},
			&cli.StringFlag{
				Name:    "user",
				Usage:   "connect as this user",
				EnvVars: []string{"PLUGIN_USER", "ANSIBLE_USER"},
			},
			&cli.StringFlag{
				Name:    "connection",
				Usage:   "connection type to use",
				EnvVars: []string{"PLUGIN_CONNECTION"},
			},
			&cli.IntFlag{
				Name:    "timeout",
				Usage:   "override the connection timeout in seconds",
				EnvVars: []string{"PLUGIN_TIMEOUT"},
			},
			&cli.StringFlag{
				Name:    "ssh-common-args",
				Usage:   "specify common arguments to pass to sftp/scp/ssh",
				EnvVars: []string{"PLUGIN_SSH_COMMON_ARGS"},
			},
			&cli.StringFlag{
				Name:    "sftp-extra-args",
				Usage:   "specify extra arguments to pass to sftp only",
				EnvVars: []string{"PLUGIN_SFTP_EXTRA_ARGS"},
			},
			&cli.StringFlag{
				Name:    "scp-extra-args",
				Usage:   "specify extra arguments to pass to scp only",
				EnvVars: []string{"PLUGIN_SCP_EXTRA_ARGS"},
			},
			&cli.StringFlag{
				Name:    "ssh-extra-args",
				Usage:   "specify extra arguments to pass to ssh only",
				EnvVars: []string{"PLUGIN_SSH_EXTRA_ARGS"},
			},
			&cli.BoolFlag{
				Name:    "become",
				Usage:   "run operations with become",
				EnvVars: []string{"PLUGIN_BECOME"},
			},
			&cli.StringFlag{
				Name:    "become-method",
				Usage:   "privilege escalation method to use",
				EnvVars: []string{"PLUGIN_BECOME_METHOD", "ANSIBLE_BECOME_METHOD"},
			},
			&cli.StringFlag{
				Name:    "become-user",
				Usage:   "run operations as this user",
				EnvVars: []string{"PLUGIN_BECOME_USER", "ANSIBLE_BECOME_USER"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Config: Config{
			Requirements:  c.String("requirements"),
			Galaxy:        c.String("galaxy"),
			Inventories:   c.StringSlice("inventory"),
			Playbooks:     c.StringSlice("playbook"),
			Limit:         c.String("limit"),
			SkipTags:      c.String("skip-tags"),
			StartAtTask:   c.String("start-at-task"),
			Tags:          c.String("tags"),
			ExtraVars:     c.StringSlice("extra-vars"),
			ModulePath:    c.StringSlice("module-path"),
			GalaxyForce:   c.Bool("galaxy-force"),
			Check:         c.Bool("check"),
			Diff:          c.Bool("diff"),
			FlushCache:    c.Bool("flush-cache"),
			ForceHandlers: c.Bool("force-handlers"),
			ListHosts:     c.Bool("list-hosts"),
			ListTags:      c.Bool("list-tags"),
			ListTasks:     c.Bool("list-tasks"),
			SyntaxCheck:   c.Bool("syntax-check"),
			Forks:         c.Int("forks"),
			VaultID:       c.String("vault-id"),
			VaultPassword: c.String("vault-password"),
			Verbose:       c.Int("verbose"),
			PrivateKey:    c.String("private-key"),
			User:          c.String("user"),
			Connection:    c.String("connection"),
			Timeout:       c.Int("timeout"),
			SSHCommonArgs: c.String("ssh-common-args"),
			SFTPExtraArgs: c.String("sftp-extra-args"),
			SCPExtraArgs:  c.String("scp-extra-args"),
			SSHExtraArgs:  c.String("ssh-extra-args"),
			Become:        c.Bool("become"),
			BecomeMethod:  c.String("become-method"),
			BecomeUser:    c.String("become-user"),
		},
	}

	if len(plugin.Config.Playbooks) == 0 {
		return errors.New("you must provide a playbook")
	}

	if len(plugin.Config.Inventories) == 0 {
		return errors.New("you must provide an inventory")
	}

	return plugin.Exec()
}
