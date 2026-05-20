# Set these to the desired values
ARTIFACT_ID=k8s-component-operator
VERSION=1.14.0
## Image URL to use all building/pushing image targets
IMAGE=cloudogu/${ARTIFACT_ID}:${VERSION}
GOTAG=1.26.0
MAKEFILES_VERSION=10.9.0
LINT_VERSION=v2.9.0
MOCKERY_VERSION=v2.53.6

ADDITIONAL_CLEAN=dist-clean

K8S_RUN_PRE_TARGETS = setup-etcd-port-forward helm-repo-config-local
PRE_COMPILE = generate-deepcopy
K8S_COMPONENT_SOURCE_VALUES = ${HELM_SOURCE_DIR}/values.yaml
K8S_COMPONENT_TARGET_VALUES = ${HELM_TARGET_DIR}/values.yaml
HELM_PRE_GENERATE_TARGETS = helm-values-update-image-version
HELM_POST_GENERATE_TARGETS = helm-values-replace-image-repo template-stage template-image-pull-policy template-log-level
IMAGE_IMPORT_TARGET=image-import
CHECK_VAR_TARGETS=check-all-vars

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
build-boot: helm-apply kill-operator-pod ## Builds a new version of the dogu and deploys it into the K8s-EcoSystem.

.PHONY: helm-values-update-image-version
helm-values-update-image-version: $(BINARY_YQ)
	@echo "Updating the image version in source value.yaml to ${VERSION}..."
	@$(BINARY_YQ) -i e ".manager.image.tag = \"${VERSION}\"" ${K8S_COMPONENT_SOURCE_VALUES}

.PHONY: helm-values-replace-image-repo
helm-values-replace-image-repo: $(BINARY_YQ)
	@if [[ ${STAGE} == "development" ]]; then \
      		echo "Setting dev image repo in target value.yaml!" ;\
    		$(BINARY_YQ) -i e ".manager.image.registry=\"$(shell echo '${IMAGE_DEV}' | sed 's/\([^\/]*\)\/\(.*\)/\1/')\"" ${K8S_COMPONENT_TARGET_VALUES} ;\
    		$(BINARY_YQ) -i e ".manager.image.repository=\"$(shell echo '${IMAGE_DEV}' | sed 's/\([^\/]*\)\/\(.*\)/\2/')\"" ${K8S_COMPONENT_TARGET_VALUES} ;\
    	fi

##@ Deployment

.PHONY: setup-etcd-port-forward
setup-etcd-port-forward:
	kubectl -n ${NAMESPACE} port-forward etcd-0 4001:2379 &

.PHONY: template-stage
template-stage: $(BINARY_YQ)
	@if [[ ${STAGE} == "development" ]]; then \
  		echo "Setting STAGE env in deployment to ${STAGE}!" ;\
		$(BINARY_YQ) -i e ".manager.env.stage=\"${STAGE}\"" ${K8S_COMPONENT_TARGET_VALUES} ;\
	fi

.PHONY: template-log-level
template-log-level: $(BINARY_YQ)
	@echo "Setting LOG_LEVEL env in deployment to ${LOG_LEVEL}!"
	@$(BINARY_YQ) -i e ".manager.env.logLevel=\"${LOG_LEVEL}\"" ${K8S_COMPONENT_TARGET_VALUES}

.PHONY: template-image-pull-policy
template-image-pull-policy: $(BINARY_YQ)
	@if [[ ${STAGE} == "development" ]]; then \
  		echo "Setting PULL POLICY to always!" ;\
		$(BINARY_YQ) -i e ".manager.imagePullPolicy=\"Always\"" ${K8S_COMPONENT_TARGET_VALUES} ;\
	fi

.PHONY: kill-operator-pod
kill-operator-pod:
	@echo "Restarting k8s-dogu-operator!"
	@kubectl -n ${NAMESPACE} delete pods -l 'app.kubernetes.io/name=${ARTIFACT_ID}'

##@ Helm-Repo-Config
.PHONY: helm-repo-config
helm-repo-config: ## Creates a configMap and a secret for the helm repo connection from env var HELM_REPO_ENDPOINT and either HELM_REPO_USERNAME & HELM_REPO_PASSWORD or HELM_AUTH_BASE64.
	@kubectl create configmap component-operator-helm-repository --from-literal=endpoint=${HELM_REPO_ENDPOINT} --from-literal=schema=oci --from-literal=plainHttp=${HELM_REPO_PLAIN_HTTP}
	@if [ -z ${HELM_AUTH_BASE64} ]; then \
	  	echo "Using fields HELM_REPO_USERNAME & HELM_REPO_PASSWORD to create secret!" ;\
		kubectl create secret generic component-operator-helm-registry --from-literal=config.json='{"auths": {"${HELM_REPO_ENDPOINT}": {"auth": "$(shell printf "%s:%s" "${HELM_REPO_USERNAME}" "${HELM_REPO_PASSWORD}" | base64 -w0)"}}}' ;\
	else \
		echo "Using field HELM_AUTH_BASE64 to create secret!" ;\
		kubectl create secret generic component-operator-helm-registry --from-literal=config.json='{"auths": {"${HELM_REPO_ENDPOINT}": {"auth": "${HELM_AUTH_BASE64}"}}}' ;\
	fi

.PHONY: helm-repo-config-local
helm-repo-config-local: ## Creates a configMap and a local config.json for the helm repo connection from env var HELM_REPO_ENDPOINT and either HELM_REPO_USERNAME & HELM_REPO_PASSWORD or HELM_AUTH_BASE64.
	@kubectl create configmap component-operator-helm-repository --from-literal=endpoint=${HELM_REPO_ENDPOINT} --from-literal=schema=oci --from-literal=plainHttp=${HELM_REPO_PLAIN_HTTP}
	@mkdir -p tmp/.helmregistry
	@if [ -z ${HELM_AUTH_BASE64} ]; then \
	  	echo "Using fields HELM_REPO_USERNAME & HELM_REPO_PASSWORD to create config.json!" ;\
		echo '{"auths": {"${HELM_REPO_ENDPOINT}": {"auth": "$(shell printf "%s:%s" "${HELM_REPO_USERNAME}" "${HELM_REPO_PASSWORD}" | base64 -w0)"}}}' > tmp/.helmregistry/config.json ;\
	else \
		echo "Using field HELM_AUTH_BASE64 to create config.json!" ;\
		echo '{"auths": {"${HELM_REPO_ENDPOINT}": {"auth": "${HELM_AUTH_BASE64}"}}}' > tmp/.helmregistry/config.json ;\
	fi

##@ Debug

.PHONY: print-debug-info
print-debug-info: ## Generates info and the list of environment variables required to start the operator in debug mode.
	@echo "The target generates a list of env variables required to start the operator in debug mode. These can be pasted directly into the 'go build' run configuration in IntelliJ to run and debug the operator on-demand."
	@echo "STAGE=$(STAGE);LOG_LEVEL=$(LOG_LEVEL);KUBECONFIG=$(KUBECONFIG);NAMESPACE=$(NAMESPACE)"

## Overwrite dev apply targets when ecosystem-core is enabled. Set ECOSYSTEM_CORE_DISABLED to true to enable legacy behaviour.
ifneq ($(ECOSYSTEM_CORE_DISABLED),true)

.PHONY: overwrite-dev-version
overwrite-dev-version: ##
	@echo "adding dev tag to image"
	$(eval IMAGE_DEV_VERSION=$(IMAGE_DEV):$(COMPONENT_DEV_VERSION))

.PHONY: helm-apply
helm-apply: overwrite-dev-version check-k8s-namespace-env-var image-import helm-chart-import ## Generates the component operator image, pushes it to the registry and then pulls the ecosystem-core chart to locally update the component operator version in the ecosystem-core helm chart (values.yaml and Chart.yaml)
	@echo "Pulling ecosystem-core chart for modification"
	@rm -rf ${K8S_RESOURCE_TEMP_FOLDER}/tmp
	@mkdir -p ${K8S_RESOURCE_TEMP_FOLDER}/tmp/ecosystem-core
	@${BINARY_HELM} pull oci://registry.cloudogu.com/k8s/ecosystem-core:$$(helm history ecosystem-core -n ${NAMESPACE} -o json | jq -r '.[-1].app_version') --untar --untardir ${K8S_RESOURCE_TEMP_FOLDER}/tmp
	@echo "Modifying Chart.yaml..."
	@${BINARY_YQ} -i '(.dependencies[] | select(.name == "k8s-component-operator") | .version) = "${COMPONENT_DEV_VERSION}"' ${K8S_RESOURCE_TEMP_FOLDER}/tmp/ecosystem-core/Chart.yaml
	@${BINARY_YQ} -i '(.dependencies[] | select(.name == "k8s-component-operator") | .repository) = "oci://registry.cloudogu.com/testing/k8s"' ${K8S_RESOURCE_TEMP_FOLDER}/tmp/ecosystem-core/Chart.yaml
	@${BINARY_HELM} dependency update ${K8S_RESOURCE_TEMP_FOLDER}/tmp/ecosystem-core
	@echo "Apply modified ecosystem-core helm chart"
	@${BINARY_HELM} --kube-context="${KUBE_CONTEXT_NAME}" upgrade -i ecosystem-core ${K8S_RESOURCE_TEMP_FOLDER}/tmp/ecosystem-core --namespace ${NAMESPACE} --reuse-values --set k8s-component-operator.manager.image.tag=${COMPONENT_DEV_VERSION} --set k8s-component-operator.manager.image.registry=registry.cloudogu.com --set k8s-component-operator.manager.image.repository=testing/$(ARTIFACT_ID)/$(GIT_BRANCH)

.PHONY: component-apply
component-apply: ## component-apply cannot be used with ecosystem-core enabled
	@echo "component-apply cannot be used with the component-operator when ecosystem-core is enabled. Use helm-apply instead."

.PHONY: component-delete
component-delete: ## component-delete cannot be used with ecosystem-core enabled
	@echo "component-delete cannot be used with the component-operator."

.PHONY: helm-delete
helm-delete: ## helm-delete cannot be used with ecosystem-core enabled
	@echo "helm-delete cannot be used with the component-operator."

endif
