General Usage
=============

Workflow
--------

General Overview
^^^^^^^^^^^^^^^^

The `gsc` tool (Git OSC) is used to bind a package source with a Git
repository. Typically, Open Build System (OBS) is acting like a
version control system, but it has a number of limitations for modern
development. One of the missing parts `gsc` is solving where package
releases and its sources (do not mix with the installed software
sources) can be tracked in a single Git repository.

Git OSC is tracking package versions in a separate branches, assisting
with the changelog maintenance per a package version.

How Does it Work
^^^^^^^^^^^^^^^^

Package, which is going to be linked to the Git repository must have a
special file in the package: `_git_pkg_repo`. That file is normally
generated automatically during package linking. It is an XML file with
two fields in it:

1. Package Git repository URL
2. Name of the release branch (typically `release-version`,
   e.g. `release-1.2`).

If for some reasons you want to create it manually, it looks like so::

  <?xml version="1.0" encoding="UTF-8"?>
  <git>
    <url>github.com:johnsmith/my-package</url>
    <branch>release-1.0</branch>
  </git>


Understanding the Workflow
^^^^^^^^^^^^^^^^^^^^^^^^^^

Currently supported workflow assumes each time package needs to be
checked out into a **private branch**, all modification done there,
and from the private branch must be performed submit request for
further steps.

It is important to understand that Git repository is a master of the
package sources, while package in the OBS is just its reflection.


Clone Package
^^^^^^^^^^^^^

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
