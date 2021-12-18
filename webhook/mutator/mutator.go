package mutator

import (
	"context"
	"fmt"

	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CuriesyncInjector injects the necessary volumes, mounts, init and sidecar containers into a pod
// necessary for synchronizing Curiefense configuration.
type CuriesyncInjector interface {
	Inject(ctx context.Context, pod *corev1.Pod) error
}

// Mutator implements the mutation webhook.
type Mutator struct {
	CuriesyncInjector CuriesyncInjector
}

func (m Mutator) Mutate(ctx context.Context, ar *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhmutating.MutatorResult, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return &kwhmutating.MutatorResult{}, nil
	}

	err := m.CuriesyncInjector.Inject(ctx, pod)
	if err != nil {
		return nil, fmt.Errorf("failed to inject curiesync: %w", err)
	}

	return &kwhmutating.MutatorResult{MutatedObject: pod}, nil
}
