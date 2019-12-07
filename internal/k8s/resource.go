package k8s

import "k8s.io/apimachinery/pkg/api/resource"

// Helper function to use a provided quantity if set, or defer to a fallback.
func resourceWithFallback(provided, fallback string) (resource.Quantity, error) {
	if provided != "" {
		return resource.ParseQuantity(provided)
	}

	return resource.ParseQuantity(fallback)
}
