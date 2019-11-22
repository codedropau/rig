package run

import (
	"github.com/alecthomas/kingpin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/nickschuch/rig/internal/compose"
	composeconfig "github.com/nickschuch/rig/internal/compose/config"
	"github.com/nickschuch/rig/internal/config"
	"github.com/nickschuch/rig/internal/k8s"
)

type command struct {
	Project string
	Master string
	Kubecfg string
	Namespace string
	Name string
	Repository string
	Tag string
	Domains []string
}

func (cmd *command) run(c *kingpin.ParseContext) error {
	cfg := config.Config{
		Dockerfiles: []string{
			"docker-compose.yml",
		},
		Services: []string{
			"nginx",
			"php-fpm",
			"mysql-default",
		},
	}

	dc, err := composeconfig.Load(cfg.Dockerfiles)
	if err != nil {
		return err
	}

	config, err := clientcmd.BuildConfigFromFlags(cmd.Master, cmd.Kubecfg)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	project, err := compose.Project(cmd.Project)
	if err != nil {
		return err
	}

	params := k8s.Params{
		Project: project,

		Namespace: cmd.Namespace,
		Name: cmd.Name,

		Repository: cmd.Repository,
		Tag: cmd.Tag,

		Domains: cmd.Domains,

		Config: cfg,
		Compose: dc,
	}

	return k8s.Apply(clientset, params)
}

// Command which snapshots a Docker Compose stack.
func Command(app *kingpin.Application) {
	c := new(command)

	cmd := app.Command("run", "Takes a snapshot of the existing Docker Compose stack.").Action(c.run)
	cmd.Flag("master", "Tag to apply to all images when performing a snapshot.").StringVar(&c.Master)
	cmd.Flag("kubecfg", "Tag to apply to all images when performing a snapshot.").Envar("KUBECONFIG").StringVar(&c.Kubecfg)
	cmd.Flag("project", "Tag to apply to all images when performing a snapshot.").Envar("RIG_PROJECT").StringVar(&c.Project)
	cmd.Flag("namespace", "Tag to apply to all images when performing a snapshot.").Required().Envar("RIG_NAMESPACE").StringVar(&c.Namespace)
	cmd.Flag("repository", "Tag to apply to all images when performing a snapshot.").Required().Envar("RIG_REPOSITORY").StringVar(&c.Repository)
	cmd.Arg("name", "Tag to apply to all images when performing a snapshot.").Required().StringVar(&c.Name)
	cmd.Arg("tag", "Tag to apply to all images when performing a snapshot.").Required().StringVar(&c.Tag)
	cmd.Arg("domain", "Domain which this environment will be exposed.").Required().StringsVar(&c.Domains)
}