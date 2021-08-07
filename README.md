# datad

`datad` is a system service that periodically collects data about its host and reports the data to a remote service. `datad` runs as a systemd service unit and will collect data on a schedule determined by configuration values.

## Collection

Data is collected through "agents". A collection agent is an executable file found in the directory `/usr/libexec/datad-agents/`. When `datad` begins a collection, it executes each agent. An agent's standard output is piped into a buffer by `datad`. This buffer is then scanned for any redactions listed in the configuration structure. Finally, the buffer is added to a map of agent output results before being submitted to the collection service.

## Submission

`datad` can be configured to submit collected data to a collection service. This is done by specifying a base URL in the `[submission]` section of `/etc/datad/config.toml`. If the submission server requires authentication, the keys `auth-type`, `username`, `password`, `cert-file`, and `key-file` may be required, depending on the configuration of the submission server. The agent output is aggregated into a table mapping an agent name to the output buffer. This data is then compressed and transmitted via an HTTP POST request to the submission server.

## Configuration

`datad` reads its main configuration from the file `/etc/datad/config.toml`. Additional configuration files may be added to `/etc/datad/config.d`, and are loaded in lexicographic order into the main configuration structure at runtime. Thus, override a value in a later file will replace any values specified previously.

### File Access Denial

Access to a file may be denied by adding an entry to the `denylist` array.

```toml
denylist = [
	"/etc/pki/.*",
]
```

Entries in the `files` key are regular expressions matching a file's absolute path. If a file is matched by an expression in this array, it is omitted entirely from any collected data.

### File Content Redaction

File content may be redacted on a line-by-line basis by adding an entry to the `redactions` array.

```toml
redactions = [
	".*lan$",
]
```

When data is collected, output is searched, line by line, using the provided regular expression pattern. Any line found to match the pattern is omitted.

