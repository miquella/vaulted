Markdown Conversion
===================

To simplify documentation writing, our man pages are first written in Markdown
and then converted to roff format using the `md2man` Ruby gem.

If you modify a Markdown document in this directory you should commit the
regenerated roff version in the same commit. To regenerate the roff version,
use this command:

```sh
md2man-roff DOCFILE.1.md > DOCFILE.1
```

But first, you may have to install the `md2man` gem: `gem install md2man`
