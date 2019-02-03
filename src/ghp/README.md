https://github.com/CurtisLusmore/ghp



### A simple web server for serving static GitHub Pages locally, to test before deploying.

This can be useful compared to browsing local HTML files from your browser when you use absolute paths in links, such as /about, /js/app.js, /css/style.css, etc., which won't resolve correctly in the context of your filesystem.

It is also handy compared to something like python -m http.server which doesn't support dropping the file extension, e.g. /about rather than /about.html.

When requesting any path ($path), ghp will do the following (all file operations are relative to the root commandline flag):

Check whether $path points to a file, if so serve that file
Check whether $path points to a directory, if so serve $path/index.html
Check whether $path.html points to a file, if so serve that file
Check whether 404.html is a file, if so serve that file as a 404
Serve a 404
If any of the above results in serving a Markdown file (extension .md), render the contents as HTML by using the GitHub Markdown API.
Getting It
From source (requires installing Go):

$ git clone https://github.com/CurtisLusmore/ghp
$ cd ghp
$ go build ghp.go
With Go Get (requires installing Go):

$ go get github.com/CurtisLusmore/ghp
Pre-compiled binaries: Check the latest Releases

Usage
$ ghp -help
Usage of ghp:
  -port int
        The port to serve over (default 8080)
  -root string
        The root directory to serve files from (your GitHub Pages repo) (default ".")
$ ghp -root MyGitHubPages
Notes
This tool currently does not support building Jekyll-based GitHub Pages. If you use Jekyll-based GitHub Pages, please see Setting up your GitHub Pages site locally with Jekyll.
As this tool exposes your filesystem to your network, you should be careful using this on untrusted networks.
This tool will send the contents of Markdown files (extension .md) to https://api.github.com/markdown. Make sure not to run this in a directory with sensitive content in Markdown files - this is mostly intended for use on files you are likely to push to GitHub Pages anyway.
Todo
Confirm response headers match live GitHub Pages