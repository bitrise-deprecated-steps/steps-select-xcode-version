# Select Xcode version

**DEPRECATED**: There's no need to use this Step on the new [bitrise.io](https://www.bitrise.io)
Stacks.

This step is basically a no-op [on the new Stacks](http://blog.bitrise.io/2016/01/20/here-comes-the-new-stack.html),
and should not be used anymore.


## Description

The new Select Xcode version step.

Selects / activates the specified Xcode version.

For available "version channel" identifiers and exact version numbers
see the [Virtual Machine pre-installed tools](http://devcenter.bitrise.io/docs/virtual-machines-updates)
article on our [DevCenter](http://devcenter.bitrise.io).

Can be run directly with the [bitrise CLI](https://github.com/bitrise-io/bitrise),
just `git clone` this repository, `cd` into it's folder in your Terminal/Command Line
and call `bitrise run test`.

*Check the `bitrise.yml` file for required inputs which have to be
added to your `.bitrise.secrets.yml` file!*

## Technical notes

This step operates based on `version-map` files.
A `version-map` is like a symlink file, contains nothing but an absolute
path to the related `Xcode.app`.

Example:

`/Applications/Xcodes/version-map--stable` might contain the path: `/Applications/Xcodes/Xcode-stable.app`

and

`/Applications/Xcodes/version-map-xcode-6` might contain the path: `/Applications/Xcodes/Xcode-6.app`

The version map files **have to be** located in `/Applications/Xcodes`
but the actual Xcode version it points to can be located anywhere.
