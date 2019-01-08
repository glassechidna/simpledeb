# simpledeb

`simpledeb` aims to be the simplest way to create an apt repo from a collection
of `.deb` files. It can be run from any platform and has no dependencies. Typical
usage is:

```
# this generates signer.key and signer.pub in local dir
simpledeb key --name "My Name" --email "me@example.com"

# this adds the new pub key as trusted by apt
cat signer.pub | apt-key add -

# this generates the apt metadata hierarchy based on the passed-in debs
simpledeb build input/*.deb

tree .
# .
# ├── input
# │   ├── example_v0.0.5-next_linux_386.deb
# │   └── example_v0.0.5-next_linux_amd64.deb
# ├── repo
# │   └── dists
# │       └── stable
# │           ├── InRelease
# │           ├── Release
# │           ├── Release.gpg
# │           └── main
# │               ├── binary-amd64
# │               │   ├── Packages
# │               │   ├── Packages.gz
# │               │   └── example_v0.0.5-next_linux_amd64.deb
# │               └── binary-i386
# │                   ├── Packages
# │                   ├── Packages.gz
# │                   └── example_v0.0.5-next_linux_386.deb
# ├── signer.key
# └── signer.pub
# 
# 7 directories, 13 files

# now if you served the contents of `repo` at example.com, you could add this to your 
# sources.list:
echo 'deb http://example.com/ stable main' >> /etc/apt/sources.list 
```

Note that simpledeb is "additive". When you run it a second time, it won't _remove_
any files from the existing repo, it will only ever add to the collection.

# Acknowledgements

This project is a fork of [`esell/deb-simple`](https://github.com/esell/deb-simple),
who did all the hard work. Any useful is thanks to them, any bugs are probably my
own work.
