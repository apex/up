import io
import json
import queue
import subprocess
import threading

proc = subprocess.Popen(["./main"], stdin=subprocess.PIPE, stdout=subprocess.PIPE, universal_newlines=True, bufsize=1)

lock = threading.Lock()
requests = queue.Queue(10)

# a hash of id -> (input_slot, event, output_slot)
responses = dict()

def send_to_proxy(proxy_stdin):
    id = 0
    while True:
        id = id + 1
        item = requests.get()
        if item is None:
            break
        cmd = item[0]

        # remove cmd from item, so it can be garbage collected faster
        item[0] = None
        cmd['id'] = id

        with lock:
            responses[id] = item

        proxy_stdin.write(json.dumps(cmd) + '\n')
        requests.task_done()

def read_from_proxy(proxy_stdout):
    while True:
        line = proxy_stdout.readline()
        if len(line) == 0:
            raise Exception('EOF from proxy')
        
        # even though there will be newlines, json.loads will ignore the whitespace
        try:
            resp = json.loads(line)
        except json.decoder.JSONDecodeError:
            # ignore lines that can't be parsed as json, someone is misuing stdout
            continue

        # If the response doesn't have an id, we can't do anything - assume misused stdout
        if 'id' not in resp:
            continue

        id  = resp['id']
        item = None
        with lock:
            if id in responses:
                item = responses[id]
                del responses[id]

        if item is not None:
            # set response into output_slot
            item[2] = resp
            # notify waiting thread that a response is ready
            item[1].set()

threading.Thread(target=send_to_proxy, name='send_to_proxy', args=(proc.stdin,)).start()
threading.Thread(target=read_from_proxy, name='read_from_proxy', args=(proc.stdout,)).start()

def handle(event, context):
    if context.identity is not None:
        identity = {
            'cognitoIdentityId': context.identity.cognito_identity_id,
            'cognitoIdentityPoolId': context.identity.cognito_identity_pool_id,
        }

    cmd = {
        # set id to 0 for now, the 'send_to_proxy' thread will assign it before serializing
        'id': 0,
        'event': event,
        'context': {
            'awsRequestId': context.aws_request_id,
            'functionName': context.function_name,
            'functionVersion': context.function_version,
            'logGroupName': context.log_group_name,
            'logStreamName': context.log_stream_name,
            'memoryLimitInMB': context.memory_limit_in_mb,
            'clientContext': context.client_context,
            'identity': identity,
            'invokedFunctionArn': context.invoked_function_arn,
        },
    }

    event = threading.Event()

    # Keep one slot of the list for the response
    item = [cmd, event, None]

    requests.put(item)
    event.wait()
    resp = item[2]

    if 'error' in resp:
        raise Exception(error)

    return resp['value']
