
var child = require('child_process');
var byline = require('./byline');

/**
 * Debug env var.
 */

const debug = process.env.DEBUG_SHIM;

/**
 * A map of string(id) to callback function, used for when
 * many concurrent requests are outstanding.
 */

var callbacks = {};

/**
 * The last id attached to a request / callback pair
 */

var lastId = (Date.now() / 1000) | 0;

/**
 * nextId generates ids which will only be repeated every 2^52 times being generated
 */

function nextId(){
  // Prevent bugs where integer precision wraps around on floating point numbers
  // (usually around 52-53 bits)
  var id = (lastId + 1) | 0;
  if (id === lastId) {
    id = 1;
  }
  lastId = id;
  return String(id);
}

/**
 * Child process for binary I/O.
 */

var proc = child.spawn('./main', { stdio: ['pipe', 'pipe', process.stderr] });

proc.on('error', function(err){
  console.error('[shim] error: %s', err);
  process.exit(1);
})

proc.on('exit', function(code, signal){
  console.error('[shim] exit: code=%s signal=%s', code, signal);
  process.exit(1);
})

/**
 * Newline-delimited JSON stdout.
 */

var out = byline(proc.stdout)

out.on('data', function(line){
  if (debug) console.log('[shim] parsing: `%s`', line)

  var msg;
  try {
    msg = JSON.parse(line);
  } catch (err) {
    console.log('[shim] unexpected non-json line: `%s`', line);
    return
  }

  if (typeof msg.id !== 'string') {
    console.log('[shim] unexpected line - do not use stdout: `%s`', line);
    return
  }

  const c = callbacks[msg.id];
  delete callbacks[msg.id];

  if (!c) {
    if (debug) console.log('[shim] unexpected duplicate response: `%s`', line)
    return
  }

  c(msg.error, msg.value);
});


/**
 * Handle events.
 */
exports.handle = function(event, ctx, cb) {
  ctx.callbackWaitsForEmptyEventLoop = false;

  const id = nextId();
  callbacks[id] = cb;

  proc.stdin.write(JSON.stringify({
    "id": id,
    "event": event,
    "context": ctx
  })+'\n');
}
