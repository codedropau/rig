package list

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/codedropau/rig/internal/k8s"
	"github.com/gosuri/uitable"
	corev1 "k8s.io/api/core/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type command struct {
	Master    string
	Kubecfg   string
	Namespace string
}

func (cmd *command) run(c *kingpin.ParseContext) error {
	config, err := clientcmd.BuildConfigFromFlags(cmd.Master, cmd.Kubecfg)
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	podList, err := clientset.CoreV1().Pods(cmd.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list Pods: %w", err)
	}

	ingressList, err := clientset.NetworkingV1beta1().Ingresses(cmd.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list Pods: %w", err)
	}

	table := uitable.New()

	table.AddRow("NAME", "VERSION", "DOMAINS", "SERVICES")

	for _, pod := range podList.Items {
		var (
			managed string
			version string
		)

		if val, ok := pod.ObjectMeta.Labels[k8s.LabelManagedBy]; ok {
			managed = val
		}

		if val, ok := pod.ObjectMeta.Labels[k8s.LabelVersion]; ok {
			version = val
		}

		if managed != k8s.ManagedBy {
			continue
		}

		table.AddRow(pod.ObjectMeta.Name, version, getDomains(ingressList, pod.ObjectMeta.Name), getServices(pod))
	}

	fmt.Println(table)

	return nil
}

// Command which lists all the environments
func Command(app *kingpin.Application) {
	c := new(command)

	cmd := app.Command("list", "Lists all the environments which are currently running.").Action(c.run)

	cmd.Flag("master", "Tag to apply to all images when performing a snapshot.").StringVar(&c.Master)
	cmd.Flag("kubecfg", "Tag to apply to all images when performing a snapshot.").Envar("KUBECONFIG").StringVar(&c.Kubecfg)

	cmd.Flag("namespace", "Namespace which the environments reside.").Required().Envar("RIG_NAMESPACE").StringVar(&c.Namespace)
}

func getDomains(list *networkingv1beta1.IngressList, name string) []string {
	var domains []string

	for _, item := range list.Items {
		if item.ObjectMeta.Name == name {
			for _, rule := range item.Spec.Rules {
				domains = append(domains, rule.Host)
			}
		}
	}

	return domains
}

func getServices(pod corev1.Pod) []string {
	var services []string

	for _, container := range pod.Spec.Containers {
		services = append(services, container.Name)
	}

	return services
}
