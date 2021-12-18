package mutation

import (
	"context"

	corev1 "k8s.io/api/core/v1"

	"github.com/sagikazarmark/kube-curiesync-injector/curiefense/curiesync"
)

const curieconfVolumeName = "curieconf"

// nolint: gosec
const (
	BucketLinkAnnotation  = "curiesync.curiefense.sagikazarmark.dev/bucket-link"
	SecretNameAnnotation  = "curiesync.curiefense.sagikazarmark.dev/secret-name"
	InitRunModeAnnotation = "curiesync.curiefense.sagikazarmark.dev/init-run-mode"
)

// CuriesyncInjector injects the necessary volumes, mounts, init and sidecar containers into a pod
// necessary for synchronizing Curiefense configuration.
type CuriesyncInjector struct {
	// CuriesyncImage is the container image to use for init and sidecar containers.
	CuriesyncImage string

	// BucketLink is the default bucket link used for synchronizing configuration.
	BucketLink string
}

func (i CuriesyncInjector) Inject(_ context.Context, pod *corev1.Pod) error {
	// Empty dir volume for the configuration
	pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
		Name: curieconfVolumeName,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	})

	bucketLink := getBucketLink(pod, i.BucketLink)
	secretName := getSecretName(pod)

	// Add init container to pod for fetching the initial configuration
	// (be it the original bootstrap config OR the real one)
	pod.Spec.InitContainers = append(pod.Spec.InitContainers, createContainer(i.CuriesyncImage, bucketLink, secretName, getInitRunMode(pod)))

	// Append volume mount to existing containers
	for i := range pod.Spec.Containers {
		pod.Spec.Containers[i].VolumeMounts = append(pod.Spec.Containers[i].VolumeMounts, corev1.VolumeMount{
			Name:      curieconfVolumeName,
			MountPath: "/config",
		})
	}

	// Add the curiesync sidecar container
	pod.Spec.Containers = append(pod.Spec.Containers, createContainer(i.CuriesyncImage, bucketLink, secretName, curiesync.PeriodicSync.String()))

	return nil
}

func createContainer(image string, bucketLink string, secretName string, runMode string) corev1.Container {
	env := []corev1.EnvVar{
		{
			Name:  "RUN_MODE",
			Value: runMode,
		},
	}

	if bucketLink != "" {
		env = append(env, corev1.EnvVar{
			Name:  "CURIE_BUCKET_LINK",
			Value: bucketLink,
		})
	}

	var envFrom []corev1.EnvFromSource
	if secretName != "" {
		envFrom = append(envFrom, corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: secretName,
				},
			},
		},
		)
	}

	return corev1.Container{
		Name:    "curiesync",
		Image:   image,
		Env:     env,
		EnvFrom: envFrom,
		VolumeMounts: []corev1.VolumeMount{
			{
				Name:      curieconfVolumeName,
				MountPath: "/config",
			},
		},
	}
}

func getBucketLink(pod *corev1.Pod, def string) string {
	if len(pod.Annotations) == 0 {
		return def
	}

	bucketLink, ok := pod.Annotations[BucketLinkAnnotation]
	if !ok || bucketLink == "" {
		return def
	}

	return bucketLink
}

func getSecretName(pod *corev1.Pod) string {
	if len(pod.Annotations) == 0 {
		return ""
	}

	return pod.Annotations[SecretNameAnnotation]
}

func getInitRunMode(pod *corev1.Pod) string {
	if len(pod.Annotations) == 0 {
		return curiesync.CopyBootstrap.String()
	}

	switch pod.Annotations[InitRunModeAnnotation] {
	case "sync-once":
		return curiesync.SyncOnce.String()

	case "bootstrap":
		return curiesync.CopyBootstrap.String()

	default:
		return curiesync.CopyBootstrap.String()
	}
}
