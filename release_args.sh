#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

# this function will be sourced from release.sh and be called from release_functions.sh
update_versions_modify_files() {
  newReleaseVersion="${1}"
  valuesYAML=k8s/helm/values.yaml
  componentPatchTplYAML=k8s/helm/component-patch-tpl.yaml
  sonarProjectProperties=sonar-project.properties

  ./.bin/yq -i ".manager.image.tag = \"${newReleaseVersion}\"" "${valuesYAML}"
  ./.bin/yq -i ".values.images.componentOperator |= sub(\":(([0-9]+)\.([0-9]+)\.([0-9]+)((?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))|(?:\+[0-9A-Za-z-]+))?)\", \":${newReleaseVersion}\")" "${componentPatchTplYAML}"
  sed -i "s/^sonar.projectVersion=.*/sonar.projectVersion=${newReleaseVersion}/" "${sonarProjectProperties}"
}

update_versions_stage_modified_files() {
  valuesYAML=k8s/helm/values.yaml
  componentPatchTplYAML=k8s/helm/component-patch-tpl.yaml
  sonarProjectProperties=sonar-project.properties

  git add "${valuesYAML}" "${componentPatchTplYAML}" "${sonarProjectProperties}"
}
