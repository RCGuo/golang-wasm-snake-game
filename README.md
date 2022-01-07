# Golang WebAssembly Snake Game

A snake game written in Golang and run in the browser by cross-compiling it to WebAssembly.

I rewrite a snake game written in Javascript to Go WASM based on [*snake-game*](https://github.com/RodionChachura/snake-game.git) developed by RodionChachura to get more familiar with Go WASM.

Demo: https://rcguo.github.io/golang-wasm-snake-game/

## Usage

`make build_wasm` Use Golang to build wasm binary file.

`make serve` To serve local http server on port 8080 (http://localhost:8080/)

> Windows user can use [Chocolatey](https://chocolatey.org/install) to install Make (choco install make) before running makefiles in Windows.

## Misc
The `wasm_exec.js` in the site directory is copy from Go 1.17.5. You can copy from your Go source to replace it:
```
$ cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./site
```

## Reference
* [Go 1.11: WebAssembly for the gophers](https://medium.zenika.com/go-1-11-webassembly-for-the-gophers-ae4bb8b1ee03)
* [Snake Game with JavaScript](https://geekrodion.medium.com/snake-game-with-javascript-10e0ad9edb52)
  * https://github.com/RodionChachura/snake-game.git
* [Go, WebAssembly, HTTP requests and Promises](https://withblue.ink/2020/10/03/go-webassembly-http-requests-and-promises.html)
* [Using Go in the Browser via Web Assembly](https://ian-says.com/articles/golang-in-the-browser-with-web-assembly/)
* https://github.com/olivewind/go-webassembly-canvas