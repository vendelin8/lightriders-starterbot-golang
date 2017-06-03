# Builds the bot for uploading

## Input arguments
* -o, --output-file: the location of the output file (default: bot.zip)
* -i, --input-dir  : the location of the source, set only if not found automatically
* -k, --keep-temp  : don't delete temp folder at the end, but writes it's location to the console

## Steps
* copies bot to a temporary directory
* removes debug.go, and all function calls to it
* zips it to a given location

## Considered, but NOT included yet
* Adding third party packages, since deploying is in an offline environment
