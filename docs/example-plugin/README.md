# Plugin Example
> :)

## How do I build a plugin?
You can use the included script, or build it manually:
```bash
# change `-X main.osTarget=windows` to reflect your build target
go build -buildmode=plugin -ldflags="-X main.osTarget=windows" cmd/command/main.go
```

## How do I install my plugin?
all you have to do is move it into your jc/plugins folder:
```bash
# if you dont have a folder, make one, the paths are listed below.
# we also check /etc/jc/plugins on unix instead of Program Files.

# user plugin
mv output.dll ~/.jc/plugins

# system plugin
mv output.dll "C:\Program Files\jc\plugins\output.dll"
```