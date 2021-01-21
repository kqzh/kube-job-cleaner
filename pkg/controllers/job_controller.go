/*
Copyright 2021 kqzh.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"kube-job-cleaner/pkg/elastic"
	"kube-job-cleaner/pkg/utils"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type JobConfig struct {
	TTLSecondsAfterFinished int

	WithConfigMap bool
}

// JobReconciler reconciles a Job object
type JobReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	EsConfig elasticsearch.Config

	Config JobConfig
}

// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs/status,verbs=get;update;patch

func (r *JobReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("job", req.NamespacedName)

	job := &batchv1.Job{}

	err := r.Client.Get(context.TODO(), req.NamespacedName, job)
	if errors.IsNotFound(err) {
		r.Log.Info("job has been deleted")

		return ctrl.Result{}, r.CleanJobESlogs(req.Name)
	}
	if err != nil {
		return ctrl.Result{}, err
	}
	if utils.IsJobFinished(job) {
		go func() {
			// waiting for ttl finished
			time.Sleep(time.Second * time.Duration(r.Config.TTLSecondsAfterFinished))

			r.Delete(context.TODO(), job, client.PropagationPolicy(metav1.DeletePropagationBackground))

			if r.Config.WithConfigMap {
				configMaps := utils.GetConfigMaps(job)
				for _, configMap := range configMaps {
					r.Delete(context.TODO(), configMap)
				}
				r.Log.Info("delete configmaps successfully")
			}
		}()
	}

	return ctrl.Result{}, nil
}

func (r *JobReconciler) CleanJobESlogs(name string) error {
	if r.EsConfig.CloudID != "" {
		esClient, err := elastic.New(r.EsConfig)
		if err != nil {
			return err
		}
		if err := esClient.DeleteLogs(name); err != nil {
			return err
		}
		r.Log.Info("clean elasticsearch logs successfully")
	}
	return nil
}

func (r *JobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.Job{}).
		Complete(r)
}
