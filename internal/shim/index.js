
var child = require('child_process')
var byline = require('./byline')

/**
 * Callback for the request.
 */

var callback

/**
 * Child process for binary I/O.
 */

var proc = child.spawn('./main', { stdio: ['pipe', 'pipe', process.stderr] })

proc.on('error', function(err){
  console.error('[shim] error: %s', err)
  process.exit(1)
})

proc.on('exit', function(code, signal){
  console.error('[shim] exit: code=%s signal=%s', code, signal)
  process.exit(1)
})

/**
 * Newline-delimited JSON stdout.
 */

var out = byline(proc.stdout)

out.on('data', function(line){
  if (process.env.DEBUG_SHIM) console.log('[shim] parsing: `%s`', line)
  var msg = JSON.parse(line)
  callback(msg.error, msg.value)
})

/**
 * Handle events.
 */

exports.handle = function(event, ctx, cb) {
  callback = cb
  ctx.callbackWaitsForEmptyEventLoop = false

  proc.stdin.write(JSON.stringify({
    "event": event,
    "context": ctx
  })+'\n');
}
