# lightriders-starterbot-golang
golang starter bot for Riddles.io's Light Riders AI challenge

## Getting started:
```
go get github.com/vendelin8/lightriders-starterbot-golang/...
```

Create a shell/batch file with content similar to this:
```
go install github.com/vendelin8/lightriders-starterbot-golang/lightRiders-starterBot-go && java -jar game-wrapper-*.jar "$(cat wrapper-commands.json)"
```

From now on, I assume you have the $GOPATH/bin in your PATH. If you don't, just call
```
export PATH=$PATH:$GOPATH/bin
```
For mac and linux, or
```
SET PATH=%PATH%;%GOPATH%\bin
```
for windows.

Add the bot to the ```wrapper-commands.json``` config like this:
```
...
"command": "lightRiders-starterBot-go"
...
```

## Replayer
<img alt="Screenshot of the replayer" style="float: right;" src="http://vendelin.square7.ch/l/sc.png" title="Replayer running in the console"/>
Call the following:

```
go install github.com/vendelin8/lightriders-starterbot-golang/replayer && replayer
```
You can add a parameter to the replayer, which is the file to replay, otherwise it will use the last one.

## Debugging
You can use the replayer to check how the bot works. It can track certain variables that will be updated in the replayer.
For adding new variables, do the following:
* Set the names of the variables in the utils/replay.go file at ```DEBUG_VARS = []string{...}```
* Update these variables while running in the lightRiders-starterBot-go/debug.go file at
```
func saveMovesToReplay() {
...
rm := utils.ReplayMove{ownBot.LastMove, oppBot.LastMove,
    []string{...}} //write debug variable values here in the same order
```
Otherwise you can use logs, called with ```logI``` for info or ```logE``` for error.
These functions call the corresponding methods of the [log15 package](http://gopkg.in/inconshreveable/log15.v2).

## Building
The uploaded package needs to build offline, so it must contain everything. This part is responsible for
combining the different sources, replacing references, and zipping the output.

For packaging the bot for uploading, call
```
go install github.com/vendelin8/lightriders-starterbot-golang/builder && builder
```

For command line arguments and more dtails check the README files of the submodules in the directories.

Have fun.

## Other Go starter bot
[Another Go bot here](https://github.com/royerk/GoLightRiders-StarterBot) to start with.
