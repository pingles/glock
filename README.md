# Glock
Go-iplemented Global Lock.

## Rationale
[Flock](http://man7.org/linux/man-pages/man1/flock.1.html) implements a filesystem lock useful when running [Cron](http://man7.org/linux/man-pages/man8/cron.8.html) if there's a chance multiple instances of a job may be executing.

Glock mimics this using ZooKeeper to provide a global lock.

## Building

Glock dependencies were managed using [Govendor](https://github.com/kardianos/govendor). You can specify `GO15VENDOREXPERIMENT=1` to pull dependencies from the `./vendor` directory.

## License

BSD 3-clause. Please see `LICENSE`.
