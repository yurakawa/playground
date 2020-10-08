#!/bin/bash

set -euxo pipefail

function increment_tag() {
    if [[ "$1" =~ ([a-z]+)-v([0-9]+).([0-9]+).([0-9]+) ]]; then
      prefix=${BASH_REMATCH[1]}
      major_version_number=${BASH_REMATCH[2]}
      minor_version_number=${BASH_REMATCH[3]}
      patch_version_number=${BASH_REMATCH[4]}

      if [ "${prefix}" = "release" ]; then
        minor_version_number=$(( minor_version_number+1 ))
        patch_version_number=0
      elif [ "${prefix}" = "hotfix" ]; then
        patch_version_number=$(( patch_version_number+1 ))
      fi

      incremented_version_number=v${major_version_number}.${minor_version_number}.${patch_version_number}
      new_tag=${prefix}-${incremented_version_number}
      echo "${new_tag}"
    else
      echo 'there may be a problem with the latest tag'
      exit 1
    fi
}

if $# -ne 1 ; then
  echo '1 argument required, e.g. release, hotfix'
  exit 1
fi

branch_type=$1

if [[ $branch_type =~ !(release|hotfix) ]]; then
  echo 'validation error, Please pass hotfix or release'
  exit 1
fi

version=$(git describe --tags --abbrev=0)
incremented_tag=$(increment_tag "${branch_type}-${version}")

if test "${GITHUB_TOKEN:-}" = ""; then
    echo "GITHUB_TOKEN must be specified"
    exit 1
fi

if ! test "$(command -v gh)"; then
    echo 'install gh command first.'
    exit 1
fi

cd "$(mktemp -d)"

service_name="playground"
git clone https://github.com/yurakawa/"${service_name}".git
cd "${service_name}"


if [ "${branch_type}" = "release" ]; then
  origin="origin/develop"
elif [ "${branch_type}" = "hotfix" ]; then
  origin="origin/master"
fi

git checkout -b "$incremented_tag" $origin
git push origin "$incremented_tag"

echo "${GITHUB_TOKEN}" > /tmp/gtoken
gh auth login --with-token < /tmp/gtoken
gh pr create -f -B master -H "${incremented_tag}" --title "${incremented_tag}"
rm -f /tmp/gtoken
