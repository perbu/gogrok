# Gogrok code analytics


## Introduction

The goal of Gogrok is to provide insight into how different Go modules interface with each other.  The primary 
question it is trying to answer is "what code imports this package?". 

![Gogrok](gogrok.png)

## How it works

Gogrok has a folder which contains all the Go modules that it is analyzing. For ease of use these should be git 
checkouts.

Gogrok works by parsing the go.mod files to map inter-module dependencies. It then proceeds to analyze the source code
to see how the different packages relate to each other.

Gogrok has no persistent storage. It parses all the code every time it starts up. Because the Go parser is so fast,
this has not been a problem for me. But if you find this problematic it should be fairly simple to serialize the 
parsed data to disk.

### How to update the git repos

If you have all the repos in a file called `repos.txt` you can GNU Parallel the following command to update all the repos: 
```bash
xargs -I {} -P 5 git -C {} pull < repos.txt
```
One repo on each line in `repos.txt`. 

This will go a lot faster if you have the following in your `~.ssh/config` file:
```
Host github.com
  ControlMaster auto
  ControlPath ~/.ssh/sockets/%r@%h-%p
  ControlPersist 600
```
This will allow you to reuse the same SSH connection for all the git interactions.

If you want to do it one repo at a time you can use the following command:
```bash
for repo in $(cat repos.txt); do echo $repo; git -C $repo pull; done
```

## Interface

The interface is a SPA that is served by the Go server. It uses HTMX to load the different fragments.

### Front page - initial state.

The front page has some navigation links at the top. There are three links. 
"local modules", "external dependencies", and "about". "local modules" is the default page.

Below there is a main fragment where the content is loaded. Initially this fragment contains the local 
modules listing.

For each module it states the module name/path,  which is clickable. When you click on the module 
name/path it will replace the main fragment with the module fragment.

Below the name/path lists the number of packages, the number of files, the number of lines of code in the module.

There are also two buttons to view dependencies and reverse dependencies. When shown the dependencies are listed in a
vertical list, each module is clickable. Clicking on the dependency will load the module fragment for the dependency.

### Local Module fragment

The module page gets invoked when you click on a module.
Below the module name it lists the followings stats for the module:
The number of packages, the number of files, the number of lines of code in the module.

Below is a horizontal list of the packages in the module. 

Clicking on a package will load the package fragment below but will not change the module fragment.

The package fragment will add another horizontal list of the files in the package. If you click on a file it will load
the file fragment below the package fragment. The file fragment will show the file.

### External dependencies fragment

This fragment lists all the external dependencies. 
It lists the module name/path. Below each module is a bullet list of the packages that import this module.