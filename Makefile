# Set these to the desired values
ARTIFACT_ID=k8s-component-operator
VERSION=0.6.0
## Image URL to use all building/pushing image targets
IMAGE=cloudogu/${ARTIFACT_ID}:${VERSION}
GOTAG?=1.21
MAKEFILES_VERSION=9.0.0
LINT_VERSION?=v1.52.1

ADDITIONAL_CLEAN=dist-clean

K8S_RUN_PRE_TARGETS=setup-etcd-port-forward
PRE_COMPILE=generate-deepcopy
HELM_PRE_APPLY_TARGETS=template-stage template-image-pull-policy template-log-level

K8S_COMPONENT_SOURCE_VALUES = ${HELM_SOURCE_DIR}/values.yaml
K8S_COMPONENT_TARGET_VALUES = ${HELM_TARGET_DIR}/values.yaml
CRD_SOURCE = ${HELM_CRD_SOURCE_DIR}/templates/k8s.cloudogu.com_components.yaml
CRD_POST_MANIFEST_TARGETS = crd-add-labels crd-copy-for-go-embedding
HELM_PRE_GENERATE_TARGETS = helm-values-update-image-version
HELM_POST_GENERATE_TARGETS = helm-values-replace-image-repo

include build/make/variables.mk
include build/make/self-update.mk
include build/make/dependencies-gomod.mk
include build/make/build.mk
include build/make/test-common.mk
include build/make/test-unit.mk
include build/make/static-analysis.mk
include build/make/clean.mk
include build/make/digital-signature.mk
include build/make/mocks.mk
include build/make/k8s-controller.mk

.PHONY: build-boot
build-boot: crd-helm-apply helm-apply kill-operator-pod ## Builds a new version of the dogu and deploys it into the K8s-EcoSystem.

.PHONY: crd-add-labels
crd-add-labels: $(BINARY_YQ)
	@echo "Adding labels to CRD..."
	@$(BINARY_YQ) -i e ".metadata.labels.app = \"ces\"" ${CRD_SOURCE}
	@$(BINARY_YQ) -i e ".metadata.labels.\"app.kubernetes.io/name\" = \"${ARTIFACT_ID}\"" ${CRD_SOURCE}

.PHONY: crd-copy-for-go-embedding
crd-copy-for-go-embedding:
	@echo "Copy CRD to pkg/api/v1/"
	@cp ${CRD_SOURCE} pkg/api/v1/

.PHONY: helm-values-update-image-version
helm-values-update-image-version: $(BINARY_YQ)
	@echo "Updating the image version in source value.yaml to ${VERSION}..."
	@$(BINARY_YQ) -i e ".controllerManager.image.tag = \"${VERSION}\"" ${K8S_COMPONENT_SOURCE_VALUES}

.PHONY: helm-values-replace-image-repo
helm-values-replace-image-repo: $(BINARY_YQ)
	@if [[ ${STAGE} == "development" ]]; then \
      		echo "Setting dev image repo in target value.yaml!" ;\
    		$(BINARY_YQ) -i e ".controllerManager.image.repository=\"${IMAGE_DEV}\"" ${K8S_COMPONENT_TARGET_VALUES} ;\
    	fi

##@ Deployment

.PHONY: setup-etcd-port-forward
setup-etcd-port-forward:
	kubectl -n ${NAMESPACE} port-forward etcd-0 4001:2379 &

.PHONY: template-stage
template-stage: $(BINARY_YQ)
	@if [[ ${STAGE} == "development" ]]; then \
  		echo "Setting STAGE env in deployment to ${STAGE}!" ;\
		$(BINARY_YQ) -i e ".controllerManager.env.stage=\"${STAGE}\"" ${K8S_COMPONENT_TARGET_VALUES} ;\
	fi

.PHONY: template-log-level
template-log-level: $(BINARY_YQ)
	@echo "Setting LOG_LEVEL env in deployment to ${LOG_LEVEL}!"
	@$(BINARY_YQ) -i e ".controllerManager.env.logLevel=\"${LOG_LEVEL}\"" ${K8S_COMPONENT_TARGET_VALUES}

.PHONY: template-image-pull-policy
template-image-pull-policy: $(BINARY_YQ)
	@if [[ ${STAGE} == "development" ]]; then \
  		echo "Setting PULL POLICY to always!" ;\
		$(BINARY_YQ) -i e ".controllerManager.imagePullPolicy=\"Always\"" ${K8S_COMPONENT_TARGET_VALUES} ;\
	fi

.PHONY: kill-operator-pod
kill-operator-pod:
	@echo "Restarting k8s-dogu-operator!"
	@kubectl -n ${NAMESPACE} delete pods -l 'app.kubernetes.io/name=${ARTIFACT_ID}'

##@ Helm-Repo-Config
.PHONY: helm-repo-config
helm-repo-config: ## Creates a configMap and a secret for the helm repo connection from env vars HELM_REPO_USERNAME, HELM_REPO_PASSWORD, HELM_REPO_ENDPOINT.
	@kubectl create configmap component-operator-helm-repository --from-literal=endpoint=${HELM_REPO_ENDPOINT} --from-literal=schema=oci --from-literal=plainHttp=${HELM_REPO_PLAIN_HTTP}
	@kubectl create secret generic component-operator-helm-registry --from-literal=config.json='{"auths": {"${HELM_REPO_ENDPOINT}": {"auth": "$(shell printf "%s:%s" "${HELM_REPO_USERNAME}" "${HELM_REPO_PASSWORD}" | base64 -w0)"}}}'

##@ Debug

.PHONY: print-debug-info
print-debug-info: ## Generates indo and the list of environment variables required to start the operator in debug mode.
	@echo "The target generates a list of env variables required to start the operator in debug mode. These can be pasted directly into the 'go build' run configuration in IntelliJ to run and debug the operator on-demand."
	@echo "STAGE=$(STAGE);LOG_LEVEL=$(LOG_LEVEL);KUBECONFIG=$(KUBECONFIG);NAMESPACE=$(NAMESPACE);DOGU_REGISTRY_ENDPOINT=$(DOGU_REGISTRY_ENDPOINT);DOGU_REGISTRY_USERNAME=$(DOGU_REGISTRY_USERNAME);DOGU_REGISTRY_PASSWORD=$(DOGU_REGISTRY_PASSWORD);DOCKER_REGISTRY={\"auths\":{\"$(docker_registry_server)\":{\"username\":\"$(docker_registry_username)\",\"password\":\"$(docker_registry_password)\",\"email\":\"ignore@me.com\",\"auth\":\"ignoreMe\"}}}"
