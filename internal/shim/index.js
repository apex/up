const child = require('child_process');

/**
 * Debug env var.
 */

const debug = process.env.DEBUG_SHIM;

/**
 * A map of string(id) to callback function, used for when
 * many concurrent requests are outstanding.
 */

const callbacks = new Map();

/**
 * The last id attached to a request / callback pair
 */

let lastId = (Date.now() / 1000) | 0;

/**
 * nextId generates ids which will only be repeated every 2^52 times being generated
 */

function nextId(){
  // Prevent bugs where integer precision wraps around on floating point numbers
  // (usually around 52-53 bits)
  let id = (lastId + 1) | 0;
  if (id === lastId) {
    id = 1;
  }

  lastId = id;
  return String(id);
}

/**
 * handleLine is responsible for taking a line of output from the child process
 * and calling the appropiate callbacks.
 */
function handleLine(line) {
  if (debug) {
    console.log('[shim] parsing: `%s`', line);
  }

  let msg;
  try {
    msg = JSON.parse(line);
  } catch (err) {
    console.log('[shim] unexpected non-json line: `%s`', line);
    return;
  }

  if (typeof msg.id !== 'string') {
    console.log('[shim] unexpected line - do not use stdout: `%s`', line);
    return;
  }

  const c = callbacks.get(msg.id);
  callbacks.delete(msg.id);

  if (!c) {
    if (debug) {
      console.log('[shim] unexpected duplicate response: `%s`', line);
    }

    return;
  }

  c(msg.error, msg.value);
}

/**
 * Child process for binary I/O.
 */

const proc = child.spawn('./main', { stdio: ['pipe', 'pipe', process.stderr] });

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

// Chunks holds onto partial chunks received in the absense of a newline.
// invariant: an array of Buffer objects, all of which do not have any newline characters
let chunks = [];
const NEWLINE = '\n'.charCodeAt(0);

// Find successive newlines in this chunk, and pass them along to `handleChunk`
function handleChunk(chunk) {
  // since this current chunk can have multple lines inside of it
  // keep track of how much of the current chunk we've consumed
  let chunkPos = 0;
  for (;;) {
    // Find the first newline in the current, in the part of the current chunk we have not
    // looked yet.
    const newlinePos = chunk.indexOf(NEWLINE, chunkPos);

    // We were not able to find any more newline characters in this chunk,
    // save the remaineder in `chunks` for later processing
    if (newlinePos === -1) {
      chunks.push(chunk.slice(chunkPos));
      break;
    }

    // We have found an end of a whole line, the beginning of the line will be the combination
    // of all Buffers currently buffered in the `chunks` array (if any)
    const start = chunk.slice(chunkPos, newlinePos);

    chunks.push(start);
    const line = Buffer.concat(chunks);
    chunks = [];

    // increase the chunk position, to skip over the last line we just found
    chunkPos = newlinePos + 1;
    handleLine(line)
  }
}

const out = proc.stdout;

out.on('readable', () => {
  for (;;) {
    const chunk = out.read();
    if (chunk === null) {
      break;
    }

    // Pump all data chunks into chunk handler
    handleChunk(chunk);
  }
});

/**
 * Handle events.
 */
exports.handle = function(event, ctx, cb) {
  ctx.callbackWaitsForEmptyEventLoop = false;

  const id = nextId();
  callbacks.set(id, cb);

  proc.stdin.write(JSON.stringify({
    "id": id,
    "event": event,
    "context": ctx
  })+'\n');
}
