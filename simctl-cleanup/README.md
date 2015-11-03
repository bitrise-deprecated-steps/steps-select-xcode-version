# simctl-cleanup

Tool to remove duplicate entries of `xcrun simctl list`.

By default this tool will only print the duplicate entries,
you can delete these by specifying the `-delete` flag.

If you want to list all available Simulators you can
call it with the `-all` flag, and if you want to only print the
Simulator IDs (so you can pipe it to another command
like `xcrun simctl delete`) you can use the `-id-only` flag.
