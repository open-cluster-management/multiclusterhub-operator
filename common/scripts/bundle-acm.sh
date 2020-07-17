#!/bin/bash
# Tested on Mac only


function extractBundleFromImage {
    id=$(docker create $1 sh)
    mkdir build/temp-bundle-$2
    docker cp $id:"/"  - > build/temp-bundle-$2/local-tar-file.tgz
    tar -xzvf build/temp-bundle-$2/local-tar-file.tgz -C build/temp-bundle-$2
    docker rm $id
    
    rm -rf bundles/$2
    mkdir bundles/$2
    cp -R build/temp-bundle-$2/manifests bundles/$2
    cp -R build/temp-bundle-$2/metadata bundles/$2
    
    rm -rf build/temp-bundle-$2
}

startBundle="quay.io/open-cluster-management/acm-operator-bundle:$1"
updateBundle="quay.io/open-cluster-management/acm-operator-bundle:$2"
startVerion=$3
updateVersion=$4

registry="quay.io/rhibmcollab"
indexImage="quay.io/rhibmcollab/multiclusterhub-operator:$updateVersion-index"

# Pull and tag images
docker pull $startBundle
docker pull $updateBundle
docker tag $startBundle $registry/acm-operator-bundle:$startVerion
docker tag $updateBundle $registry/acm-operator-bundle:$updateVersion
startBundle="$registry/acm-operator-bundle:$startVerion"
updateBundle="$registry/acm-operator-bundle:$updateVersion"

# Extract Contents of Images
extractBundleFromImage $startBundle $startVerion
extractBundleFromImage $updateBundle $updateVersion

# Add 'Replaces' to Update CSV
yq w -i \
    bundles/$updateVersion/manifests/advanced-cluster-management.v$updateVersion.clusterserviceversion.yaml \
    "spec.replaces" "advanced-cluster-management.v$startVerion"

# Switch all channels to latest
yq w -i \
    bundles/$startVerion/metadata/annotations.yaml \
    "annotations.[operators.operatorframework.io.bundle.channels.v1]" "latest"
yq w -i \
    bundles/$startVerion/metadata/annotations.yaml \
    "annotations.[operators.operatorframework.io.bundle.channel.default.v1]" "latest"
yq w -i \
    bundles/$updateVersion/metadata/annotations.yaml \
    "annotations.[operators.operatorframework.io.bundle.channels.v1]" "latest"
yq w -i \
    bundles/$updateVersion/metadata/annotations.yaml \
    "annotations.[operators.operatorframework.io.bundle.channel.default.v1]" "latest"

# Generate and Build Bundles
opm alpha bundle generate \
--directory bundles/$startVerion/manifests \
--package advanced-cluster-management \
--channels latest \
--default latest \

docker build -f bundle.Dockerfile -t $startBundle .

opm alpha bundle generate \
--directory bundles/$updateVersion/manifests \
--package advanced-cluster-management \
--channels latest \
--default latest

docker build -f bundle.Dockerfile -t $updateBundle .

rm bundle.Dockerfile

# Push new bundles
docker push $startBundle
docker push $updateBundle

# Generate and push index image of bundles
opm index add \
--bundles $startBundle,$updateBundle \
--tag $indexImage \
-c docker

docker push $indexImage

# Update Kustomize with new tag if necessary
yq w -i build/index-install/downstream/kustomization.yaml 'images[0].newTag' "$updateVersion-index"
