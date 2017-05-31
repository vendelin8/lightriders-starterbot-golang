# lightriders-starterbot-golang
golang starter bot for Riddles.io's Light Riders competition

## WARNING:
A big refactor is coming soon. Until that, the code is just a proof of concept for the engine creators. Please not rely on it yet.

## Getting started:
```
go get github.com/vendelin8/lightriders-starterbot-golang
```

Create a shell/batch file with content similar to this:
```
go install github.com/vendelin8/lightriders-starterbot-golang/lightRiders-starterBot-go && java -jar game-wrapper-*.jar "$(cat wrapper-commands.json)"
```


I assume you have the $GOPATH/bin in your PATH. Add the bot to the ```wrapper-commands.json``` config like this:
```
...
"command": "lightRiders-starterBot-go"
...
```

If you add the bot to the second place, it will automatically save replays to a directory called ```replays```.

If you want to replay, call
```
go install github.com/vendelin8/lightriders-starterbot-golang/replayer && replayer
```
You can add a parameter to the replayer, which is the file to replay, otherwise it will use the last one.

Have fun.

## Other Go starter bot
[Another Go bot here](https://github.com/royerk/GoLightRiders-StarterBot) to start with.
