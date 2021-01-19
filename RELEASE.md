# Releasing

Releases of this repo target [semver][semver] tags pushed by repo admins.
To request a release, please ping an admin on [#olm-dev][slack-olm-dev]
or [#operator-sdk-dev][slack-osdk-dev] Kubernetes Slack channels, or
post to the [operator-framework group][of-ggroup].

## Tags

As per semver, all releases containing new features must map to a major or minor version increase.
Patch releases must only contain fixes to features released in a prior release.

## Process

In your local shell (assuming you have repo admin privileges):

```sh
export PREVIOUS_RELEASE_TAG=$(git describe --tags --abbrev=0)
export RELEASE_TAG="vX.Y.Z"
git checkout master
git pull master
git fetch --all
git tag $RELEASE_TAG
# Assuming the 'upstream' remote points to the operator-framework repo.
git push upstream refs/tags/$RELEASE_TAG
```

Then create release notes while still on the `master` branch:

```sh
while read -r line; do echo $line | awk '{f = $1; $1 = ""; print "-"$0; }'; done <<< $(git log $PREVIOUS_RELEASE_TAG..$RELEASE_TAG --format=oneline --no-merges)
```

Copy them into the Github release [description form][release-desc-page],
select `vX.Y.Z` in the `Tag version` form, and click `Publish release`.

[semver]:https://semver.org/
[slack-olm-dev]:https://kubernetes.slack.com/messages/olm-dev
[slack-osdk-dev]:https://kubernetes.slack.com/messages/operator-sdk-dev
[of-ggroup]:https://groups.google.com/forum/#!forum/operator-framework
[release-desc-page]:https://github.com/operator-framework/api/releases/new
