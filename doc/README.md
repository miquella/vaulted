Markdown Conversion
===================

Markdown files are all generated and committed, so if you edit
any of the markdown files you will need to follow the procedure
outlined below.

Setup
-----

First, you may have to install the `md2man` gem: `gem install md2man`,
and then you may need to install the `go-bindata` dependency with:

```sh
gem install md2man
```

```sh
go get -u github.com/jteeuwen/go-bindata
```

Generating Man Pages
--------------------

To simplify documentation writing, our man pages are first written in Markdown
and then converted to roff format using the `md2man` Ruby gem.

If you modify a Markdown document in this directory you should commit the
regenerated roff version in the same commit. To regenerate the roff version,
use this command:

```sh
go generate
```
