Simple Git hook for bitrise
===========================
This simple tool is designed to send api trigger to bitrise server from a git repository hook.

Who is this for?
----------------
If you are hosting your own git server, or your tool's web hook is not yet supported by bitrise, this is a simple tool to allow you to trigger a build on bitrise whenever you push / update code in your git repo.

How to use it?
--------------
Currently there are 2 suported git hooks:
* post-update
* post-receive

The only difference right now is post-receive would send the commit hash too, where that is not available to post-update.

Compile the main.go and deploy it in the repository's hook directory, rename it to post-receive or post-update suiting your need.


