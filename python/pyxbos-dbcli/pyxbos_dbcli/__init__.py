__version__ = '0.1.0'

import asyncio
from pygments.lexers.python import PythonLexer
from prompt_toolkit import PromptSession
from prompt_toolkit.eventloop.defaults import use_asyncio_event_loop
from prompt_toolkit.patch_stdout import patch_stdout
from prompt_toolkit.styles import Style
from prompt_toolkit.lexers import PygmentsLexer

import numpy as np

session = PromptSession()

# Tell prompt_toolkit to use the asyncio event loop.
use_asyncio_event_loop()

from prompt_toolkit.completion import WordCompleter
completer = WordCompleter(["get","set","list","pending","commit"])

_pending_vars = {}
_vars = {}

def get(key, default=None):
    """
    Retrieve the value of a variable stored in this db.
    Defaults to the provided 'default' value if the key is not found
    """
    global _vars
    return _vars.get(key)

def put(key, value):
    """
    Set the value of 'key' to the result of evaluating the Python expression in
    'value'.
    """
    global _pending_vars, _vars
    _pending_vars[key] = eval(value)
    return _pending_vars[key]

def commit():
    global _pending_vars, _vars
    committed = []
    _newvars = _vars.copy()
    for k,v in _pending_vars.items():
        print(f'{k} = {v}')
        _newvars[k] = v
        committed.append(k)
    for k in committed:
        _pending_vars.pop(k)
    _pending_vars = {}
    _vars = _newvars
    return f'committed: {len(committed)}'

async def run(inp):
    global _pending_vars, _vars
    cmd, *args = inp.split(' ')
    if cmd == 'get':
        result = await session.prompt('  [key]: ', async_=True, lexer=None)
        if key in _pending_vars:
            return f'(uncommitted) {_pending_vars[key]}'
        return _vars.get(key)
    elif cmd == 'set':
        key = await session.prompt('  [key]: ', async_=True, lexer=None)
        value = await session.prompt('  [value]: ', async_=True, lexer=PygmentsLexer(PythonLexer))
        session.lexer = None
        return put(key, value)
    elif cmd == 'list':
        keys = set(list(_pending_vars.keys()))
        keys = keys.union(set(list(_vars.keys())))
        return '\n'.join(keys)
    elif cmd == 'pending':
        res = []
        for k,v in _pending_vars.items():
            res.append(f"{k} => {v} (type: {type(v)})")
        return '\n'.join(res)
    elif cmd == 'commit':
        return commit()
    elif  cmd == 'help':
        return """
Commands:

Ctl-d to exit

get: get variable
set: set variable to result of Python expression (propose)
list: list all variables
pending: list pending variables (vars you have 'set' without running 'commit')
commit: saves all proposed variables so they can be used
help: prints this
        """
    else:
        return "Command not found. Use [get, set, list, pending, commit, help]"

async def prompt_coroutine():
    print('Ctl-d to exit')
    while True:
        with patch_stdout():
            try:
                result = await session.prompt('[energise]>>> ', async_=True,
                                      completer=completer, complete_while_typing=True, lexer=None)
                print(await run(result))
            except KeyboardInterrupt:
                continue
            except EOFError:
                asyncio.get_event_loop().stop()
                break
            except Exception as e:
                print(f'Exception: {e}')

def prompt():
    """
    Open a prompt within the current asyncio event loop (make sure that is running!)
    """
    asyncio.ensure_future(prompt_coroutine())

if __name__ == '__main__':
    put('a', "123")
    put('b', '"hello"')
    commit()
    asyncio.ensure_future(prompt_coroutine())
    loop = asyncio.get_event_loop()
    loop.run_forever()
    loop.close()
