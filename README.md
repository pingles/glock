# Crony
Cron with companions.

## Usage

```
$ crony --lockPath=/some/path --schedule="* * * * *" --command="/usr/bin/env echo 'hello, world'"
```

Crony uses ZooKeeper to acquire a distributed lock, ensuring only a single process will execute a task. The task is executed according to the specified cron expression.

For example, to run `echo 'hello, world'` every second (on a single machine) you can run as follows:

`lockPath` is used to uniquely identify the task: it should be set to the same value across all machines that would attempt to execute the same command; only a single machine (from those specifying the same lock path) will execute the command.

## Rationale
Occasionally you want to run tasks that need to execute on a schedule. Like cron ;-)

Doing this on an unknown set of machines is a little more difficult if you want only one instance of the task to execute (perhaps a database backup, export etc.).

Normally one machine is the special one: tasked as the scheduler, worker etc. This could be through running a single stable instance or providing metadata (if running in a cloud environment) for a single machine. Either situation require manual intervention and control to ensure the machine is available. In [adrianco](https://twitter.com/adrianco) parlance: we must treat the host as a pet rather than cattle.

Crony solves the problem generally: allowing a command to be executed on a regular cron schedule but self-coordinating with other machines to ensure only a single task is run. It provides single-execution cron on cattle.

## Building

Crony dependencies were managed using [Govendor](https://github.com/kardianos/govendor). You can specify `GO15VENDOREXPERIMENT=1` to pull dependencies from the `./vendor` directory.

## License

BSD 3-clause. Please see `LICENSE`.
