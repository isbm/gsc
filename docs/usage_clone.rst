Cloning the Package
===================

Same as one would branch a package with `osc bco ...`, the same one
need to clone the package. For example, if project is
`home:my-project` and package is `my-package`, then the following
command would clone the package::

  gsc clone home:my-project my-package

The command above will assume that the package is already linked to
the Git repository. If there is no Git repository yet, but package
already has a long history, new *empty* Git repository must be created
before `clone` operation, and then linked on clone stage as follows::

  gsc clone home:my-project my-package git@github.com:johnsmith/my-package

In this case `gsc` assumes that the package `my-package` is not yet
linked, and it will setup that for you.

Once package is successfully cloned, it will be in a subdirectory of
your OBS home. Navigating to the package sources, you will find that
the package sources are also Git repository and is already checked out
into a temporary branch, from which you will do your changes to the
package.

Once you've done with the changes and want to perform package Submit
Request, you may consider making Pull Request to the GitHub (or Merge
Request if you use GitLab). This step is optional.

Once Submit Request is accepted, you can run `gsc merge` and it will
merge for you current temporary branch.
