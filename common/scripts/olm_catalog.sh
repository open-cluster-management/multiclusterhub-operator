#!/bin/bash

set -e

indent() {
  local INDENT="      "
  local INDENT1S="    -"
  sed -e "s/^/${INDENT}/" \
      -e "1s/^${INDENT}/${INDENT1S} /"
}

listCSV() {
  for index in ${!CSVDIRS[*]}
  do
    indent apiVersion < "$(ls "${CSVDIRS[$index]}"/*version.yaml)"
  done
}

DEPLOYDIR=${DIR:-$(cd "$(dirname "$0")"/../../deploy && pwd)}

cp "${DEPLOYDIR}"/operator.yaml "${DEPLOYDIR}"/operator.yaml.bak
if [ "$(uname)" = "Darwin" ]; then
  sed -i "" "s|multicloudhub-operator:latest|${IMAGE}|g" "${DEPLOYDIR}"/operator.yaml
else
  sed -i "s|multicloudhub-operator:latest|${IMAGE}|g" "${DEPLOYDIR}"/operator.yaml
fi
operator-sdk generate csv --csv-channel "${CSV_CHANNEL}" --csv-version "${CSV_VERSION}" >/dev/null 2>&1
cp "${DEPLOYDIR}"/operator.yaml.bak "${DEPLOYDIR}"/operator.yaml
rm -f "${DEPLOYDIR}"/operator.yaml.bak


BUILDDIR=${DIR:-$(cd "$(dirname "$0")"/../../build && pwd)}
OLMOUTPUTDIR="${BUILDDIR}"/_output/olm
mkdir -p "${OLMOUTPUTDIR}"

PKGDIR="${DEPLOYDIR}"/olm-catalog/multicloudhub-operator
CSVDIRS[0]=${DIR:-$(cd "${PKGDIR}"/"${CSV_VERSION}" && pwd)}

CRD=$(grep -v -- "---" "$(ls "${DEPLOYDIR}"/crds/*crd.yaml)" | indent)
PKG=$(indent packageName < "$(ls "${PKGDIR}"/*multicloudhub-operator.package.yaml)")
CSVFILE="${PKGDIR}"/"${CSV_VERSION}"/multicloudhub-operator.v"${CSV_VERSION}".clusterserviceversion.yaml

# remove replaces field
sed -ie '/replaces:/d' "${CSVFILE}"

# disable all namespaces supported, see https://github.com/operator-framework/operator-sdk/issues/2173 
index=$(grep -n "type: AllNamespaces" "${CSVFILE}" | cut -d ":" -f 1)
index=$((index - 1))
if [ "$(uname)" = "Darwin" ]; then
  sed -i "" "${index}s/true/false/" "${CSVFILE}"
else
  sed -i "${index}s/true/false/" "${CSVFILE}"
fi

NAME=${NAME:-multicloudhub-operator-registry}
NAMESPACE=${NAMESPACE:-multicloud-system}
DISPLAYNAME=${DISPLAYNAME:-multicloudhub-operator}

cat > "${OLMOUTPUTDIR}"/multicloudhub.resources.yaml <<EOF | sed 's/^  *$//'
# This file was autogenerated by 'common/scripts/olm_catalog.sh'
# Do not edit it manually!
---
apiVersion: operators.coreos.com/v1alpha1
kind: CatalogSource
metadata:
  name: $NAME
spec:
  configMap: $NAME
  displayName: $DISPLAYNAME
  publisher: Red Hat
  sourceType: configmap
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: $NAME
data:
  customResourceDefinitions: |-
$CRD
  clusterServiceVersions: |-
$(listCSV)
  packages: |-
$PKG
EOF

\cp -r "${PKGDIR}" "${OLMOUTPUTDIR}"
rm -rf "${DEPLOYDIR}"/olm-catalog

echo "Created ${OLMOUTPUTDIR}/olm-catalog/multicloudhub-operator"
echo "Created ${OLMOUTPUTDIR}/multicloudhub.resources.yaml"