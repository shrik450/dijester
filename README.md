# Dijester

Dijester is a simple, UNIX-y utility for generating an EPUB/Markdown "digest"
from your news sources. You can define your sources via a TOML configuration
file, and Dijester will fetch the latest articles from those sources, format
them into a digest, and save them as an EPUB or Markdown file, ready for you to
read on your device of choice.

## Goals

1. Extensible: Dijester's core architecture is designed to be easy to extend
   via writing code. Simply fork and add your own sources to fetch from
   whatever sites you use.
2. Simple: Dijester does one thing and does it well. Concerns like scheduling,
   delivery etc. are left to the user. Consider setting dijester up on a cron
   schedule, and sending the output to your Kindle via email.
