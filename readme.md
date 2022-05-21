# Dijester

Dijester is a small, self-contained utility for generating periodical digests
from online feeds.

## Usage

```shell
dijester ./dijester.toml
```

You can also

```shell
dijester --help
```

To get a list of flags and what they do.

## How it works

Dijester is a single binary that runs on Linux and macOS. When you
want to generate a digest, you simply execute dijester with a config
file (more on that later). Dijester then reads your feeds and compiles
the content into an output format of your choice, which you can then
read at your leisure.

### What it can do

1. Read any number of RSS feeds and compile entries from them.
2. Keep track of entries it's already compiled.
3. Generate a folder full of markdown files, or one compiled epub file.
4. Make a callback network request after it completes generating the
   digest.
5. Interface with a service encapsulating Mozilla's Readbility to load
   articles from feeds which do not include the full content. In the
   future, this capability _may_ be bundled into the dijester binary,
   but until then you will need to use something like
   [phpdockerio/readability-js-server](https://hub.docker.com/r/phpdockerio/readability-js-server)

### What it won't

1. Handle scheduling. You can do this yourself with cron, or use
   a more complex workflow automation system like Hugin.
2. Handle styling. If you want that, export to a bundle of files and
   use another tool to compile those into your preferred format with
   styling.
3. Handle receiving e-mails from newsletters etc. You can use an
   email-to-rss service like KillTheNewsletter for that.

## Configs

Dijester is configured using a TOML file. The details of configuration
are explained [here](./doc/config.md)
