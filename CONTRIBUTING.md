### Welcome

We welcome contributions to the OpenRadiant Project in many forms, and there's always plenty to do!

First things first, please review the OpenRadiant Project's [Code of Conduct](CONDUCT.md) before participating. It is important that we keep things civil.

### Getting help
If you are looking for something to work on, or need some expert assistance in debugging a problem or working out a fix to an issue, our community is always eager to help. We hang out on \[TODO] [Slack](https://OpenRadiant.slack.com/), IRC (#OpenRadiant on freenode.net) and the [mailing lists](http://example.com/). \[TODO] create mailing list

### Reporting bugs
If you are a user and you find a bug, please submit an [issue](https://github.com/ORG/REPO/issues). Please try to provide sufficient information for someone else to reproduce the issue. One of the project's maintainers should respond to your issue within 24 hours. If not, please bump the issue and request that it be reviewed.

### Fixing issues and working stories
Review the [issues list](https://github.ibm.com/alchemy-containers/openradiant/issues) and find something that interests you. You could also check the ["help wanted"](https://github.ibm.com/alchemy-containers/openradiant/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22) list. It is wise to start with something relatively straight forward and achievable. Usually there will be a comment in the issue that indicates whether someone has already self-assigned the issue. If no one has already taken it, then add a comment assigning the issue to yourself, eg.: ```I'll work on this issue.```. Please be considerate and rescind the offer in comments if you cannot finish in a reasonable time, or add a comment saying that you are still actively working the issue if you need a little more time.

We are using the [GitHub Flow](https://guides.github.com/introduction/flow/) process to manage code contributions. If you are unfamiliar, please review that link before proceeding.

## Setup your workstation for development

  1. Create a [fork](https://help.github.com/articles/fork-a-repo/) (if you haven't already)
  1. Clone it locally

  ```bash
  git clone https://github.com/<github_username>/containercafe.git
  ```
  1. Add the upstream repository as a remote

  ```bash
  cd containercafe
  git remote add upstream https://github.com/containercafe/containercafe.git
  ```

## To work on something, whether a new feature or a bugfix:

  1. Always update your master ([more detail here](https://help.github.com/articles/syncing-a-fork/))

  ```bash
  git checkout master
  git fetch upstream master
  git merge --ff-only upstream/master
  git push origin master
  ```

  1. Create a descriptively-named branch off of your cloned fork, called a
  *feature branch* (reference the corresponding issue id when fixing an issue)

  ```bash
  git checkout -b issue-nnnn
  ```
  
  1. Repeat the following steps until you are ready to submit for review.

    1. Work on your code updates, then commit to that branch
    locally. Regularly push your work to the same feature branch in
    your fork in github.

    1. Commit messages

      Commit messages must have a short description, no longer than 50
      characters, followed by a blank line and a longer, more
      descriptive message. It should briefly describe what you changed
      and why. Include reference(s) to issue(s) being addressed. Don't
      use ```Closes #1234``` or ```Fixes #1234``` because they cause
      the issue to be closed after merging but we want them in 'done'
      state in ZenHub for further validation.  See
      [the list of auto-close keywords of github](https://help.github.com/articles/closing-issues-via-commit-messages/)
      and avoid them because we use ZenHub.

    ```bash
    git commit
    git push origin issue-nnnn
    ```

  1. (Optionally) Squash

    Once you are ready to open a PR for your work, if you have made
    multiple commits then you may choose to squash them together into
    one commit.  If the commit from which you started working was
    `deadbeef` then you can squash using interactive rebase as
    follows.

    ```bash
    git rebase -i deadbeef
    ```

  1. Open PR (Pull Request)

    When you need feedback or help, or you think the branch is ready
    for merging, open a pull request.  Make sure you have first
    successfully built and tested your changes; explicitly state in
    your PR whether you have done so.

     _Note: if your PR does not merge cleanly or your feature branch needs to be
     updated, follow the steps below:_

        ```bash
        git checkout issue-nnnn
        git fetch upstream
        git rebase upstream/master
        git push -f origin issue-nnnn
        ```

  1. Did we mention tests? All code changes should be accompanied by new or modified tests.

  1. Continuous Integration (CI): Be sure to check [Travis](https://travis-ci.org/)
  results associated with your pull request. It provides a good indication whether
  your change is syntactically sound, if it has any conflicts or if requires
  a branch update. \[TODO] You can also check in the Slack  [#ci-status](https://containercafe.slack.com/messages/ci-status) channel for
  status of your build.

  1. Any code changes that affect documentation should be accompanied by
  corresponding changes (or additions) to the documentation and tests. This will
  ensure that if the merged PR is reversed, all traces of the change will be
  reversed as well.

After your Pull Request (PR) has been reviewed and signed off, a maintainer will
 merge it into the master branch.

## Coding guidelines
Participation in the OpenRadiant community is governed by the OpenRadiant
[Code of Conduct](CONDUCT.md)

### Becoming a maintainer
Projects or sub-projects will be lead by a set of maintainers. New projects can
designate an initial set of maintainers that will be approved by the Technical
Steering Committee when the project is first approved. The project's maintainers
will, from time-to-time, consider adding a new maintainer. An existing maintainer
will post a pull request to the [MAINTAINERS.txt](MAINTAINERS.txt) file. If a
majority of the maintainers concur in the comments, the pull request is then
merged and the individual becomes a maintainer.

### Legal stuff
We have tried to make it as easy as possible to make contributions. This applies to how we handle the legal aspects of contribution. We use the same approach&mdash;the [Developer's Certificate of Origin 1.1 (DCO)](DCO1.1.txt)&mdash;that the Linux&reg; Kernel [community](http://elinux.org/Developer_Certificate_Of_Origin) uses to manage code contributions.
We simply ask that when submitting a pull request, the developer must include a sign-off statement in the pull request description.

Here is an example Signed-off-by line, which indicates that the submitter accepts the DCO:

```
Signed-off-by: John Doe <john.doe@hisdomain.com>
```
