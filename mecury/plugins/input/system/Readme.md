# Systems monitor input plugins


### Configuration:

This is a sample configuration for the plugin.

```toml
# Read metrics about cpu usage
[[inputs.cpu]]
    ## Whether to report per-cpu stats or not
    percpu = true
    ## Whether to report total system cpu stats or not
    totalcpu = true
    ## Comment this line if you want the raw CPU time metrics
    fielddrop = ["time_*"]
    #[inputs.cpu.tagdrop]
    #    cpu = [ "cpu*" ]
    interval = "10s"
[[inputs.disk]]

[[inputs.diskio]]


[[inputs.mem]]

[[inputs.swap]]

[[inputs.net]]

[[inputs.netstat]]

[[inputs.processes]]

[[inputs.system]]
```
