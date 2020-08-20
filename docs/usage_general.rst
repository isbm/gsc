General Usage
=============

What it is About?
-----------------

The ``gsc`` tool is to make a link between the `OBS (Open Build
Service) <http://openbuildservice.org>`_  and Git repository per a
package. The Git repository is keeping all the versions of the package
source (do not mix with the software component sources!) and is an
ultimate source to all packages across the OS releases.

Typically, Open Build System (OBS) is acting like a version control
system, but it has a number of limitations for modern development. One
of the missing parts ``gsc`` is solving where package releases and its
sources (do not mix with the installed software sources) can be
tracked in a single Git repository.

Git OSC is tracking package versions in a separate branches, assisting
with the changelog maintenance per a package version.

How Does it Work
----------------

Package, which is going to be linked to the Git repository must have a
special file in the package: ``_git_pkg_repo``. That file is normally
generated automatically during package linking. It is an XML file with
two fields in it:

1. Package Git repository URL
2. Name of the release branch (typically ``release-version``,
   e.g. ``release-1.2``).

If for some reasons you want to create it manually, it looks like so::

  <?xml version="1.0" encoding="UTF-8"?>
  <git>
    <url>github.com:johnsmith/my-package</url>
    <branch>release-1.0</branch>
  </git>


Understanding the Workflow
--------------------------

The GSC is also encoding a workflow with it, where no more ``osc``
or ``git`` commands needed to be used directly anymore. For example:

* Cloning, making changes, committing to the branches at OBS and Git
  at the same time, keeping in sync everything.
* Changelog entry is pre-constructed out of Git commits.
* Submit Requests are done in a separate Git branches and then merged
  afterwards on SR being accepted.
* Tracking of all package sources (e.g. random uncovered changes to
  the .spec file is not possible).
* Each release has its own Git branch and each package can be bound to
  a specific branch.
* "Reflashing" package from the Git source (or importing one from Git
  source).
* More...

Currently supported workflow assumes each time package needs to be
checked out into a **private branch**, all modification done there,
and from the private branch must be performed submit request for
further steps.

It is important to understand that Git repository is a master of all
sources of the package, while package in the OBS is just its
reflection and can be just deleted/purged permanently, then re-created
again from the Git repo.
