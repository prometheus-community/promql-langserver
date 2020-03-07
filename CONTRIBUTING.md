# Contributing

# Internal dependencies

This repo has a dependencies that are in a marked as internal in their source repo and cannot be easily imported.

To manage these special steps are required.

For each of them there exists a respective directory in `vendored/<name>` and a couple of metadata files:

* `<name>.repo` contains the source repo of the dependency
* `<name>.version` contains the commit id that should be checked out before copying the files over in this repo
* `<name>.dirs` contains the directories that should be copied over
* `<name>.cmd` contains a `sed` command to be run on all imported files to fix import paths

To update these dependencies, set the appropriate version and run `make update_internal_packages`. Then fix every problem the update caused. This should be done every few weeks or so.