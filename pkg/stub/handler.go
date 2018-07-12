package stub

import (
	"context"

	"github.com/slintes/hello-operator/pkg/apis/hello/v1alpha1"

	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func NewHandler() sdk.Handler {
	return &Handler{}
}

type Handler struct {
	// Fill me
}

func (h *Handler) Handle(ctx context.Context, event sdk.Event) error {
	switch o := event.Object.(type) {
	case *v1alpha1.Hello:
		err := sdk.Create(newHelloPod(o))
		if err != nil && !errors.IsAlreadyExists(err) {
			logrus.Errorf("Failed to create hello pod : %v", err)
			return err
		}
	}
	return nil
}

// newHelloPod creates a hello pod
func newHelloPod(cr *v1alpha1.Hello) *corev1.Pod {
	labels := map[string]string{
		"app": "hello-" + cr.Name,
	}
	return &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "hello-" + cr.Name,
			Namespace: cr.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cr, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "Hello",
				}),
			},
			Labels: labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "messenger",
					Image:   "busybox",
					Command: []string{"/bin/sh"},
					Args:    []string{"-c", "echo " + cr.Spec.Message + "; while true; do sleep 10; done"},
				},
			},
		},
	}
}
