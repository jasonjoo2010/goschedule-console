# goschedule-console

Console tool for goschedule

## Quick Start

```shell
go build .
./goschedule-console -v
```

Then you can open `http://127.0.0.1:8000/` in your browser to use it.

## Parameters

Para | Comment
--- | ---
-v | Enable verbose logging(default to disabled)
-p | Set serving port(default to 8000)
-c | Specify config file location(default to "config.yml" in workdir, configuration will be written to it)
