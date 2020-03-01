## Go FileStorm

Golang implementation of the Storm protocol.

## Building the source

Building `storm` requires both a Go (version 1.10 or later) and a C compiler. You can install
them using your favourite package manager. Once the dependencies are installed, run

```shell
make storm
```
## Running `storm`

Going through all the possible command line flags is out of scope here (please consult our
[CLI Wiki page](https://github.com/filestorm/go-filestorm/wiki/Command-Line-Options)),
but we've enumerated a few common parameter combos to get you up to speed quickly
on how you can run your own `storm` instance.

## License

The go-filestorm library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html),
also included in our repository in the `COPYING.LESSER` file.

The go-filestorm binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also
included in our repository in the `COPYING` file.
