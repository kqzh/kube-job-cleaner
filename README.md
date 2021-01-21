## Kubernetes Job Cleaner

Kubernetes Controller to automatically delete completed Jobs and Pods and ConfigMaps and ElasticSearch Logs.

Init with kubebuilder

Some common use-case scenarios:

* Delete Jobs and their pods after their completion with ttl seconds
* Delete Jobs and their pods and pods's configmaps
* Delete Jobs and their pods and pods's logs which send to elasticSearch

| flag name       | default | usage                                      |
| --------------- |---------|------------------------------------------- |
| delete-after    | 0       | delete job and pods after specified period |
| with-configmap  | true    | delete configmaps which mount on job's pod |

| env                     | usage                     |
| ----------------------- | --------------------------|
| ElasticSearchCloudID    | Login Cloud ElasticSearch |
| ElasticSearchUsername   | -                         |
| ElasticSearchPassword   | -                         | 



