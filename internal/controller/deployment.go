package controller

import (
	"context"
	v1 "github.com/Nageshbansal/deployment-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	log "sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *DeploySetReconciler) Deployment(ctx context.Context, req ctrl.Request, ds *v1.DeploySet) (*appsv1.Deployment, error) {

	log := log.FromContext(ctx)

	replicas := ds.Spec.Replica.Count

	labels := map[string]string{
		"app.kubernetes.io/name":       "DeploySet",
		"app.kubernetes.io/instance":   ds.Name,
		"app.kubernetes.io/version":    "v1",
		"app.kubernetes.io/part-of":    "deployset-operator",
		"app.kubernetes.io/created-by": "controller-manager",
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ds.Name,
			Namespace: ds.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           ds.Spec.Container.Image,
						Name:            ds.Name,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{{
							ContainerPort: int32(ds.Spec.Container.Port),
							Name:          "deployset",
						}},
					}},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(ds, dep, r.Scheme); err != nil {
		log.Error(err, "failde to set controller owner reference")
		return nil, err
	}
	return dep, nil
}

func (r *DeploySetReconciler) DeploymentIfNotExist(ctx context.Context, req ctrl.Request, ds *v1.DeploySet) (bool, error) {

	log := log.FromContext(ctx)

	dep := &appsv1.Deployment{}

	err := r.Get(ctx, types.NamespacedName{
		Name:      ds.Name,
		Namespace: ds.Namespace,
	}, dep)

	if err != nil && errors.IsNotFound(err) {
		dep, err := r.Deployment(ctx, req, ds)
		if err != nil {
			log.Error(err, "Failed to define the new deployment in deploy set")

		}
		log.Info("Creating new deployment, Deployment.Namespace: %v, Deployment.Name: %v", dep.Namespace, dep.Name)

		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(
				err, "Failed to create new Deployment",
				"Deployment.Namespace", dep.Namespace,
				"Deployment.Name", dep.Name,
			)

			return false, err
		}

		err = r.getDeploySet(ctx, req, ds)
		if err != nil {
			log.Error(err, "Failde to fetch the deployment")
			return false, err
		}

		return true, nil
	}
	if err != nil {
		log.Error(err, "Failed to get the Deployment")
		return false, err
	}
	return false, nil

}

func (r *DeploySetReconciler) UpdateDeploymentReplica(ctx context.Context, req ctrl.Request, ds *v1.DeploySet) error {

	log := log.FromContext(ctx)

	dep := &appsv1.Deployment{}

	err := r.Get(ctx, types.NamespacedName{Name: ds.Name, Namespace: ds.Namespace}, dep)
	if err != nil {
		log.Error(err, "failed to get the Deployment")
		return err
	}

	replicas := ds.Spec.Replica.Count

	if replicas == *dep.Spec.Replicas {
		return nil
	}

	log.Info("Updating the Deployment replica")

	dep.Spec.Replicas = &replicas
	err = r.Update(ctx, dep)
	if err != nil {
		log.Error(
			err, "Failed to update Deployment",
			"Deployment.Namespace", dep.Namespace,
			"Deployment.Name", dep.Name,
		)

		err = r.getDeploySet(ctx, req, ds)
		if err != nil {
			log.Error(err, "Failed to re-fetch TDSet")
			return err
		}

		return nil
	}

	err = r.getDeploySet(ctx, req, ds)
	if err != nil {
		log.Error(err, "Failed to re-fetch TDSet")
		return err
	}
	if err != nil {
		return err
	}

	return nil
}
