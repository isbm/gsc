Change Request
==============

What is Change Request?
-----------------------

Change request in OBS terminology is a part of `Staging Workflow
<https://openbuildservice.org/help/manuals/obs-user-guide/cha.obs.best-practices.webuiusage.html#staging_how_to>`_
and is a process of making required changes to a package. Such change
affects all the resulting release and therefore must yield to a
certain rules.

Submitting Change Request
-------------------------

When you are done making changes to your package via ``gsc add``
command, then it is a time to request them to be accepted by a package
maintainer and added to the release, eventually. In order to submit a
Change Request, issue the following command::

  gsc sr

Or::

  gsc submitreq

This will bring a newly constructed changelog entry, which will
contain all previously made commits to Git. Each commit will be made
as a separate changelog entry. If needed, you should edit and reformat
that text accordingly.

After you quit changelog text editor, you might have few more text
editors displaying diff between the files. Those you should edit also
accordingly to your needs.

Once all is committed, at this point, your temporary Git branch
(``tmp-<package-name>-<version>``) will be synchronised with the Git
and Change Request will be sent to the OBS and you will see a Change
Request number at the end.

When CR Was Rejected
--------------------

In case of rejection was due to a mistakes or missing parts, then
nothing special happens, just keep working on changes by issuing ``gsc
add`` command, bringing things to the desired state and try to submit
your Change Request again.

But ss long as you think your effort no longer needed and you should
stop continue with this branch permanently, you have to clean it up
and close it gracefully::

  gsc close


When CR Was Accepted
--------------------

If your Change Request has been accepted, it also means that the
branch of your package has been accepted. There is no direct
integration between the "Accept Request" button on OBS and your Git
repository. That should be done manually. Since your branch is your
local one and wasn't sent to Git as Pull Request (because the review
of the package is happening at OBS side anyways), you should now merge
your working branch to the master or the release branch. To do so, you
just issue the following command::

  gsc mg

Or::

  gsc merge

This will merge your temporary branch to the base derived branch
(either master if development, or some release branch) without fast
forwarding and without squashing anything, so the commits are granular
and can be easily reverted.

Your merge branch will contain an automatic message::

  Merge branch tmp-mygreatpackage-1.2

You can also verify your Git log if everything is correct::

  git log

