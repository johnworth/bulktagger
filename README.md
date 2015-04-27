# bulktagger

A very basic command-line tool to tag and push multiple Docker images at once.

## Usage

    ./bulktagger --registry <registry> --pull-tag <tag> --tag <tag> --list <path-to-list>

The --list option should be a path to a file containing the names of Docker images, one on each line. Example:

    kifshare
    jex-events
    user-sessions
    user-preferences

--pull-tag is the tag of the image that will get tagged.
--tag is the new tag that will get applied to the pulled images.
--registry is the registry that should be used when pushing and pulling images.

Here's a more concrete example:

    ./bulktagger --registry discoenv --pull-tag dev --tag qa --list images.txt

That will pull the "discoenv" images contained in images.txt that have been tagged with "dev" from the Docker Hub, apply the "qa" tag to the images, and push the newly tagged qa images back up to the Docker Hub.

This assumes you've already run "docker login".
