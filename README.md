# Dijester

Dijester is a simple, UNIX-y utility for generating an EPUB/Markdown "digest"
from your news sources. You can define your sources via a TOML configuration
file, and Dijester will fetch the latest articles from those sources, format
them into a digest, and save them as an EPUB or Markdown file, ready for you to
read on your device of choice.

## Goals

1. Extensible: Dijester's core architecture is designed to be easy to extend
   via writing code. Simply fork and add your own sources to fetch from
   whatever sites you use. If you think a source is useful to others, please
   submit a pull request!
2. Simple: Dijester does one thing and does it well. Concerns like scheduling,
   delivery etc. are left to the user. Consider setting dijester up on a cron
   schedule, and sending the output to your Kindle via email.

## Usage

1. Dijester is distributed as a statically linked single executable for various
   platforms. You can download it from the
   [releases page](https://github.com/shrik450/dijester/releases). You might
   want to rename the executable to `dijester` for convenience.
2. Create a configuration file in TOML format. You can refer to the samples in
   this repository for examples, and the
   [configuration documentation](docs/configuring.md) for details on how to set
   up your sources and processors.
3. Run the executable with the path to your configuration file as an argument:

   ```bash
   ./dijester -config /path/to/config.toml
   ```

   If you want to change the output directory, you can use the `-output-dir` flag:

   ```bash
   ./dijester -config /path/to/config.toml -output-dir /path/to/output
   ```

4. The output will be saved in the specified directory as an EPUB or Markdown
   file, depending on your configuration.

## Contributing

Contributions are welcome! If you have a feature request, bug report, or
enhancement, please open an issue or submit a pull request. I reserve the right
to reject any ideas, bug reports, feature requests, or pull requests that I
consider to be out of scope for this project. If you think your idea is
important and I reject it, feel free to fork!

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file
for details.

## Acknowledgements

This project wouldn't be possible without the following libraries and tools:

1. [Shiori](https://github.com/go-shiori), which provides the excellent
   go-readability and go-epub libraries.
