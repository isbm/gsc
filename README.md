# gsc: Git to OSC Wrapper

An attempt to bring OBS and Git together without rewriting Source Server part in OBS.

## General Usage

The idea is based that the `.osc` and `.git` resides in the same place. Everthing
that is in the package directory is also going to Git. But `.git` is not going to
the package.

## How It Works

There is an OSC-magic file in the package: `_pkg_git_repo`. It contains an URL of
Git where the sources of the package are stored.

### Checkout "bco"
If no `_pkg_git_repo` file found, then the package sources considered not yet bound
to any Git repository and so a new Git repo will be initialised.
