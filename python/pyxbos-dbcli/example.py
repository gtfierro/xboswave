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
