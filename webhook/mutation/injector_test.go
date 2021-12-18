package mutation_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/sagikazarmark/kube-curiesync-injector/webhook/mutation"
)

func TestCuriesyncInjector(t *testing.T) {
	tests := []struct {
		name        string
		injector    mutation.CuriesyncInjector
		pod         *corev1.Pod
		expectedPod *corev1.Pod
	}{
		{
			"OK",
			mutation.CuriesyncInjector{
				CuriesyncImage: "curiefense/curiesync:latest",
			},
			&corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx-ingress",
							Image: "nginx-ingress:1.0.0",
						},
					},
				},
			},
			&corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{},
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:  "curiesync",
							Image: "curiefense/curiesync:latest",
							Env: []corev1.EnvVar{
								{
									Name:  "RUN_MODE",
									Value: "COPY_BOOTSTRAP",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "curieconf",
									MountPath: "/config",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "nginx-ingress",
							Image: "nginx-ingress:1.0.0",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "curieconf",
									MountPath: "/config",
								},
							},
						},
						{
							Name:  "curiesync",
							Image: "curiefense/curiesync:latest",
							Env: []corev1.EnvVar{
								{
									Name:  "RUN_MODE",
									Value: "PERIODIC_SYNC",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "curieconf",
									MountPath: "/config",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "curieconf",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
		{
			"With bucket link",
			mutation.CuriesyncInjector{
				CuriesyncImage: "curiefense/curiesync:latest",
			},
			&corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						mutation.BucketLinkAnnotation: "s3://bucket/prefix/",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx-ingress",
							Image: "nginx-ingress:1.0.0",
						},
					},
				},
			},
			&corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						mutation.BucketLinkAnnotation: "s3://bucket/prefix/",
					},
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:  "curiesync",
							Image: "curiefense/curiesync:latest",
							Env: []corev1.EnvVar{
								{
									Name:  "RUN_MODE",
									Value: "COPY_BOOTSTRAP",
								},
								{
									Name:  "CURIE_BUCKET_LINK",
									Value: "s3://bucket/prefix/",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "curieconf",
									MountPath: "/config",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "nginx-ingress",
							Image: "nginx-ingress:1.0.0",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "curieconf",
									MountPath: "/config",
								},
							},
						},
						{
							Name:  "curiesync",
							Image: "curiefense/curiesync:latest",
							Env: []corev1.EnvVar{
								{
									Name:  "RUN_MODE",
									Value: "PERIODIC_SYNC",
								},
								{
									Name:  "CURIE_BUCKET_LINK",
									Value: "s3://bucket/prefix/",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "curieconf",
									MountPath: "/config",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "curieconf",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
		{
			"With secret",
			mutation.CuriesyncInjector{
				CuriesyncImage: "curiefense/curiesync:latest",
			},
			&corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						mutation.SecretNameAnnotation: "secret",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx-ingress",
							Image: "nginx-ingress:1.0.0",
						},
					},
				},
			},
			&corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						mutation.SecretNameAnnotation: "secret",
					},
				},
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:  "curiesync",
							Image: "curiefense/curiesync:latest",
							Env: []corev1.EnvVar{
								{
									Name:  "RUN_MODE",
									Value: "COPY_BOOTSTRAP",
								},
							},
							EnvFrom: []corev1.EnvFromSource{
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "secret",
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "curieconf",
									MountPath: "/config",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:  "nginx-ingress",
							Image: "nginx-ingress:1.0.0",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "curieconf",
									MountPath: "/config",
								},
							},
						},
						{
							Name:  "curiesync",
							Image: "curiefense/curiesync:latest",
							Env: []corev1.EnvVar{
								{
									Name:  "RUN_MODE",
									Value: "PERIODIC_SYNC",
								},
							},
							EnvFrom: []corev1.EnvFromSource{
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "secret",
										},
									},
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "curieconf",
									MountPath: "/config",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "curieconf",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			err := test.injector.Inject(context.Background(), test.pod)
			require.NoError(t, err)

			assert.Equal(t, test.expectedPod, test.pod)
		})
	}
}
