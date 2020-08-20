Making Changes
==============

After package being cloned, it is now a time to start making some
changes to the package. While it is possible to make one big change to
everything and name it as one entry, it is generally a better practice
to make related changes in a smaller steps. For example, adding two
patches can be done as one change and described as::

  - Added foo.patch and bar.patch (<JIRA-FOO-ID> and <JIRA-BAR-ID>)

When one adds a patch, at least new file is required (patch itself)
and ``.spec`` file needs to be updated as well. In this case, we will
have three changes: two new files and one to the ``.spec`` file. Also
we will have one diff, covering both fixes.

But it is much better to split that into two entries::

  - Added foo.patch (<JIRA-FOO-ID>)
  - Added bar.patch (<JIRA-BAR-ID>)

To do so, each time you want to make an entry to the Change log, you
should call ``add`` command within the package directory (where
``.osc`` and ``.git`` subdirectories are located)::
    
  gsc add

This command will bring your ``$EDITOR`` with the typical OSC's
editor, ready to form a description entry about the change. You should
describe what the change is for, save and exit. You can verify your
change by issuing command::

  git log

Essentially, command ``add`` means "add a change". And it is better
that the command is called more often and more granular steps are then
resulting in the Change Log of the package.
