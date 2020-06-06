# wui

Lightweight & easy web-based GUI for your app.

Pronounced "woo-ee", it stands for "Web User interface". It's written in go, and uses a client/server model. Think of it as a super-light, cross-platform alternative to electron. All of your code is meant to run inside the context of the web-page that makes up the UI, and just trigger the server to do native-things.

## usage

```
mkdir mything
cd mything
npm init wui
npm start
```

Make a app/settings.json file to control how it operates:

```json
{
  "title": "My app",
  "icon": "/icon.png",
  "url": "/",
  "tray": {
    "icon": "/icon.png",
    "url": "/tray",
  }
}
```


### running

* Run `npm run build` to build complete zips for all supported platforms.
* Run `npm start` to run the runtime in the background, and open your web-view in a browser for development.


### API

The native API is just REST calls, but a promise-based helper library is included.

You can use it like this, inside your app, to do native things:


```html
<script type=module>
import * as WUI from '/wui.js'


/// SETTINGS

// get settings
const settings = await WUI.settings()

// merge new settings with existing (and bring app to that state)
const newSettings = await WUI.settings({
  title: 'New title',
  icon: '/newicon.png'
})


/// FILESYSTEM

// write a file
await WUI.write('test.txt', 'o hai!')

// read a file
const text = await WUI.read('test.txt')

// create a dir
await WUI.mkdir('./mydir')

// get a list of fiels in a directory
const files = await WUI.ls('./mydir')

// get detailed info about a file/directory
const info = await WUI.stat('./mydir')

// delete a directory or file
await WS.rm('./mydir')


/// ENVIRONMENT VARIABLES

// get all
const env = await WUI.env()

// get one
const PATH = await WUI.env('PATH')

// set one
const PATH = await WUI.env('COOL', '1')



/// SUB-PROCESS

// run a command
const [ out, err ] = await WS.exec('ls -al')


/// EXIT

// handle user trying to close app
WUI.on('exit', () => {
  console.log('I am trying to exit. Return (or resolve) false to stop that.')
  return false
})

// close app
WUI.exit()


</script>
```

Other than that, you can use all the regular HTML stuff, like:

* [mess with files](https://developer.mozilla.org/en-US/docs/Web/API/File/Using_files_from_web_applications)
* [open a new window](https://developer.mozilla.org/en-US/docs/Web/API/Window/open)
* [small key/value store](https://developer.mozilla.org/en-US/docs/Web/API/Web_Storage_API)
* [database](https://developer.mozilla.org/en-US/docs/Web/API/IndexedDB_API)
* dialogs: [confirm](https://developer.mozilla.org/en-US/docs/Web/API/Window/confirm), [alert](https://developer.mozilla.org/en-US/docs/Web/API/Window/alert), [prompt](https://developer.mozilla.org/en-US/docs/Web/API/Window/prompt)
* [2D](https://developer.mozilla.org/en-US/docs/Web/API/Canvas_API)
* [3D](https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API)