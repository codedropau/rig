package run

import (
	"github.com/alecthomas/kingpin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/skpr/awsutils/eks"

	composeconfig "github.com/codedropau/rig/internal/compose/config"
	"github.com/codedropau/rig/internal/config"
	"github.com/codedropau/rig/internal/k8s"
)

type command struct {
	Config string

	// Authentication used when connecting to a cluster.
	Username string
	Password string

	Master  string
	Kubecfg string

	// Metadata applied to Kubernetes objects.
	Namespace string
	Name      string

	// Information used to run the correct images.
	Repository string
	Tag        string

	// Domains which the environment will be accessible from.
	Domains []string
}

type AWS struct {
	Cluster eks.Cluster

}

func (cmd *command) run(c *kingpin.ParseContext) error {
	cfg, err := config.Load(cmd.Config)
	if err != nil {
		return err
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

	params := k8s.Params{
		Project: cfg.Project,

		// Metadata applied to Kubernetes objects.
		Namespace: cmd.Namespace,
		Name:      cmd.Name,

		// Information used to run the correct images.
		Repository: cmd.Repository,
		Tag:        cmd.Tag,

		// Domains which the environment will be accessible from.
		Domains: cmd.Domains,

		Config:  cfg,
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

	cmd.Flag("username", "Username used to authenticate with the Kubernetes cluster.").Required().Envar("RIG_USERNAME").StringVar(&c.Username)
	cmd.Flag("password", "Password used to authenticate with the Kubernetes cluster.").Required().Envar("RIG_PASSWORD").StringVar(&c.Password)

	cmd.Flag("config", "Config file to load.").Default(".rig.yml").Envar("RIG_CONFIG").StringVar(&c.Config)

	// Metadata applied to Kubernetes objects.
	cmd.Flag("namespace", "Tag to apply to all images when performing a snapshot.").Required().Envar("RIG_NAMESPACE").StringVar(&c.Namespace)
	cmd.Arg("name", "Name of the Kubernetes objects which will be provisioned.").Required().StringVar(&c.Name)

	// Information used to run the correct images.
	cmd.Flag("repository", "Tag to apply to all images when performing a snapshot.").Required().Envar("RIG_REPOSITORY").StringVar(&c.Repository)
	cmd.Arg("tag", "Tag to apply to all images when performing a snapshot.").Required().StringVar(&c.Tag)

	// Domains which the environment will be accessible from.
	cmd.Arg("domain", "Domain which this environment will be exposed.").Required().StringsVar(&c.Domains)
}