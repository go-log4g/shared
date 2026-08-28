# go-log4g-shared

`go-log4g-shared` provides a `sharedRollingFile` appender for applications using `go-log4g`.

The appender is designed for cases where **multiple processes write to the same rolling log file**.

Unlike a regular `rollingFile` appender, access to the active log file and rollover operations is coordinated using a cross-process file lock. This prevents concurrent processes from writing or rolling the shared file at the same time.

## Installation

```bash
go get github.com/go-log4g/shared
```

Import the package in the application:

```go
import _ "github.com/go-log4g/shared"
```

The package registers the `sharedRollingFile` appender with `go-log4g` automatically.

## Configuration

```yaml
appenders:
  shared:
    type: sharedRollingFile
    file: logs/shared.log
    filePattern: logs/arch/shared-%d{yyyyMMddHHmm}-%i.log.zip

    # Optional. Defaults to .lock in the directory containing the log file.
    # fileLock: logs/.lock

    # Optional. Default: true
    # append: true

    # Optional. Default: true
    # immediateFlush: true

    # Optional. Default: 8192
    # bufferSize: 8192

    policies:
      onStartupTriggeringPolicy:
        minSize: 1B

      timeBasedTriggeringPolicy:
        interval: 1
        # modulate: false

      # sizeBasedTriggeringPolicy:
      #   size: 1MB

    defaultRolloverStrategy:
      # max: 7

      delete:
        maxAge: 30d
        # maxFiles: 100
        # maxTotalSize: 5GB

    layout:
      type: pattern
      pattern: "${pattern}"
```

## Shared file locking

Every write to the shared log file is protected by a cross-process file lock.

While holding the lock, the appender:

1. reads the current state of the shared log file;
2. evaluates the rollover policy;
3. performs rollover when necessary;
4. opens the active log file;
5. appends the log data;
6. closes the file.

The lock is then released.

This ensures that the rollover decision is based on the current file state and that another process cannot modify the shared file between the rollover check and the write.

### `fileLock`

The lock file can be configured explicitly:

```yaml
fileLock: logs/.lock
```

If `fileLock` is omitted, the appender uses `.lock` in the directory containing the shared log file.

All processes writing to the same shared log must resolve to the same lock.

## Buffering

By default:

```yaml
immediateFlush: true
```

each log event reaches the shared file writer immediately.

For high-throughput logging, buffering can be enabled:

```yaml
immediateFlush: false
bufferSize: 8192
```

Each process maintains its own in-memory buffer. When a buffer is flushed, the entire buffered batch is written while holding the cross-process file lock.

This reduces the number of file opens, writes, closes, and cross-process lock acquisitions.

Because buffers are process-local, log records from different processes may appear in batches rather than being strictly interleaved by timestamp.

## Rolling

`sharedRollingFile` supports the standard rolling-file configuration provided by `go-log4g`, including:

* startup triggering;
* time-based triggering;
* size-based triggering;
* indexed file patterns using `%i`;
* date patterns using `%d{...}`;
* gzip and ZIP archives;
* rollover retention;
* deletion by age, file count, and total size.

For example:

```yaml
filePattern: logs/arch/shared-%d{yyyyMMddHHmm}-%i.log.zip
```

creates ZIP-compressed rollover files using the rollover time and archive index.

## Startup rollover

A startup policy can roll an existing shared file when a process starts:

```yaml
policies:
  onStartupTriggeringPolicy:
    minSize: 1B
```

Startup rollover is also performed while holding the shared file lock, so multiple processes starting concurrently cannot perform the operation simultaneously.

## Runtime failures

Runtime file-system failures do not terminate the application.

If a write or rollover operation fails, the error is reported through the `go-log4g` status logger. Subsequent log events can retry the operation, allowing the appender to recover from transient conditions such as temporary file access failures.

## When to use `sharedRollingFile`

Use `sharedRollingFile` when multiple processes or application instances on the same filesystem intentionally write to one log file.

For a log file owned by only one process, use the standard `rollingFile` appender instead. It avoids the additional cross-process synchronization required by `sharedRollingFile`.
