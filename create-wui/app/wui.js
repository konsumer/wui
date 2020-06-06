// This is the client-side library for WUI

window.WUI = {}

// run a remote command
const remote = (command = 'ping', params = {}, endpoint = '/_api/') => fetch(`${endpoint}${command}`, { method: 'POST', body: JSON.stringify(params) }).then(r => r.json())

/// SETTINGS

// get/set settings
window.WUI.settings = (params = {}, endpoint = '/_api/') => remote('settings', params, endpoint)

/// FILESYSTEM

// write file
// TODO: look into typed arrays for dealing with binary files
window.WUI.write = (filename = '', contents = '', endpoint = '/_api/') => remote('write', { filename, contents }, endpoint)

// read file
// TODO: look into typed arrays for dealing with binary files
window.WUI.read = (filename = '', endpoint = '/_api/') => remote('read', { filename }, endpoint).then(r => r.contents)

// make a directory
window.WUI.mkdir = (dirname = '', endpoint = '/_api/') => remote('mkdir', { dirname }, endpoint)

// list files in a directory
window.WUI.ls = (dirname = '', endpoint = '/_api/') => remote('ls', { dirname }, endpoint)

// get file/dir stats
window.WUI.stat = (filename = '', endpoint = '/_api/') => remote('stat', { filename }, endpoint)

// delete a file/directory
window.WUI.rm = (filename = '', endpoint = '/_api/') => remote('rm', { filename }, endpoint)

/// ENVIRONMENT VARIABLES

// get/set env-var(s)
window.WUI.env = (name = '', value = '', endpoint = '/_api/') => remote('env', { name, value }, endpoint)

/// SUB-PROCESS

// run a command
window.WUI.exec = (command = '', endpoint = '/_api/') => remote('exec', { command }, endpoint)

/// EXIT
window.WUI.exit = (endpoint = '/_api/') => remote('exit', {}, endpoint)

/// EVENT-HANDLERS

// TODO: need to implement this, maybe with web-sockets?
window.WUI.on = (eventName, handler) => {}
