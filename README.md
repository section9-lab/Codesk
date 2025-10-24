# README

## About

This is the official Wails React-TS template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.

### 1. build bin
```sh
cd ~/Documents/GitHub/Codesk
wails build -clean -platform darwin/universal
```
### 2. package dmg

```sh
create-dmg \
  --volname "Codesk" \
  --window-pos 200 120 \
  --window-size 400 400 \
  --icon-size 100 \
  --icon "Codesk.app" 200 190 \
  --app-drop-link 600 185 \
  --hide-extension "Codesk.app" \
  "/Codesk_v0.0.1.dmg" \
  "build/bin/"
```

Ref:
> https://github.com/winfunc/opcode
