// This is the client-side library for WUI

// run a remote command
const remote = (command = 'ping', params = {}, endpoint='/_api/') => fetch(`${endpoint}${command}`, { method: 'POST', body: JSON.stringify(params) }).then(r => r.json())


/// SETTINGS

// get/set settings
export const settings = (params = {}, endpoint='/_api/') => remote('settings', params, endpoint)


/// FILESYSTEM

// write file
// TODO: look into typed arrays for dealing with binary files
export const write = (filename = '', contents = '', endpoint='/_api/') => remote('write', { filename, contents }, endpoint)

// read file
// TODO: look into typed arrays for dealing with binary files
export const read = (filename = '', endpoint='/_api/') => remote('read', { filename }, endpoint)

// make a directory
export const mkdir = (dirname = '', endpoint='/_api/') => remote('mkdir', { dirname }, endpoint)

// list files in a directory
export const ls = (dirname = '', endpoint='/_api/') => remote('ls', { dirname }, endpoint)

// get file/dir stats
export const stat = (filename = '', endpoint='/_api/') => remote('stat', { filename }, endpoint)

// delete a file/directory
export const rm = (filename = '', endpoint='/_api/') => remote('rm', { filename }, endpoint)


/// ENVIRONMENT VARIABLES

// get/set env-var(s)
export const env = (name = '', value = '', endpoint = '/_api/') => remote('env', { name, value }, endpoint)


/// SUB-PROCESS

// run a command
export const exec = (command='', endpoint='/_api/') => remote('exec', { command }, endpoint)


/// EXIT
export const exit = (endpoint='/_api/') =>  remote('exit', {}, endpoint)


/// EVENT-HANDLERS

// TODO: need to implement this, maybe with web-sockets?
export const on = (eventName,  handler) => {}