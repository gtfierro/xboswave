# PyXBOS DB CLI

This is a convenient package for introducing an interactive prompt for getting/setting configuration values while an async program is running.

This is easiest seen with an example

```python
import pyxbos_dbcli as dbcli
import asyncio

# print out the values of a,b,c,d
# every five seconds by loading the latest
# values from dbcli
async def print_values():
    while True:
        a = dbcli.get('a')
        b = dbcli.get('b')
        c = dbcli.get('c')
        d = dbcli.get('d')
        print(f'Values:\na: {a}\nb: {b}\nc: {c}\nd: {d}\n')
        await asyncio.sleep(5) 


# make sure we run the print_values coroutine
asyncio.ensure_future(print_values())

# asynchronously start the prompt (this doesn't block)
dbcli.prompt()

# set some default values (optional)
#
# values are Python expressions that will be evaluated.
# You can use external libraries (e.g. numpy) if they have been imported
dbcli.put('a', '3')
dbcli.put('b', '[-2,1,3]')
dbcli.put('c', '"hello"')
dbcli.put('d', '{"value": 4.0}')

# need to commit values in order for them to be read by dbcli.get
dbcli.commit()

# run asyncio event loop
asyncio.get_event_loop().run_forever()
```

## Usage

Get committed values using `get(key)`. `key` is a string of the name of the configuration variable

Set new values using `put(key, value)`, where `key` is a string of the name of the configuration variable and `value` is a string containing a Python expression whose evaluation yields the value we want to store under the key.
These expressions can be strings, numbers, dictionaries, lists or even calls to external libraries like numpy (as long as numpy is imported into your program).

In order for new values to be accessed using `get`, you need to run `commit()`, which atomically replaces the old values with the new values. This way, your program will see the new values all at once, rather than bit by bit as you set them interactively.

## CLI

When you call `dbcli.prompt()`, an interactive prompt will appear on your screen. Hit Tab to get an autocomplete of available commands, or just type `help` followed by Enter to get a list of commands:

Prompt:
```
[energise]>>>
```

Entering the `help` command:

```
[energise]>>> help

Commands:

Ctl-d to exit

get: get variable
set: set variable to result of Python expression (propose)
list: list all variables
pending: list pending variables (vars you have 'set' without running 'commit')
commit: saves all proposed variables so they can be used
help: prints this
```

Just type the name of the command and it will prompt you for the correct arguments.

List current vars:

```
[energise]>>> list
a
c
b
d
```

Setting some values:

```
[energise]>>> set
  [key]: a
  [value]: 4
4
[energise]>>> set
  [key]: b
  [value]: [4,5,6]
4
```

Checking what values we have set so far (but have not committed)

```
[energise]>>> pending
a => 4 (type: <class 'int'>)
b => [4, 5, 6] (type: <class 'list'>)
```

Committing our changes

```
[energise]>>> commit
a = 4
b = [4, 5, 6]
committed: 2
```

## API

You have a programmatic API available to you. Note that the values in this package are global across all of its imports in the same process.

`get(key, default=None) -> value`: returns the value of the given variable stored in the DB. `key` is a string of the name of the configuration variable. Returned value defaults to the provided `default` value if the key is not found.

`put(key, value) -> evaluated value`:  Set new values using `put(key, value)`, where `key` is a string of the name of the configuration variable and `value` is a string containing a Python expression whose evaluation yields the value we want to store under the key.
These expressions can be strings, numbers, dictionaries, lists or even calls to external libraries like numpy (as long as numpy is imported into your program).
Returns the result of evaluating the expression in `value`.

`commit()`: Commits the current values proposed by calling `put`.

`prompt()`: Opens a prompt. Must have an asyncio event loop running
