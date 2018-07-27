majorVersion=$(cat version)
DRK_VERSION="${majorVersion}${TRAVIS_JOB_ID}"

set -e

#Create a release
releaseId=$(curl -s -H "Authorization: token $GITHUB_ACCESS_TOKEN" --data '{"tag_name": "'${DRK_VERSION}'","target_commitish": "master","name": "'${DRK_VERSION}'","body": "Release of version '${DRK_VERSION}'","draft": true,"prerelease": false}' https://api.github.com/repos/labaneilers/drake/releases | jq -r '.id')

echo "Created release $release $version"

# Upload files
pushd out
    for filename in *; do
        echo " - Uploading $filename..."
        curl -s -H "Authorization: token $GITHUB_ACCESS_TOKEN" -H "Content-Type:application/octet-stream" --data-binary @"$filename" 'https://uploads.github.com/repos/labaneilers/drake/releases/'$releaseId'/assets?name='$filename | jq -r '.url'
    done
popd

# Make it "live"
curl -s -H "Authorization: token $GITHUB_ACCESS_TOKEN" --data '{"draft": false}' 'https://api.github.com/repos/labaneilers/drake/releases/'$releaseId




# Get the status
#releaseId=$(curl -H "Authorization: token $GITHUB_ACCESS_TOKEN" https://api.github.com/repos/labaneilers/drake/releases/latest | jq -r '.id')

# Delete a release
#curl -X "DELETE" -H "Authorization: token $GITHUB_ACCESS_TOKEN" https://api.github.com/repos/labaneilers/drake/releases/12136029
