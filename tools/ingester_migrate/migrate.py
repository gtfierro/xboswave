from influxdb import InfluxDBClient

client = InfluxDBClient('localhost', 8086, '', '', 'xbos')

measurements = client.get_list_measurements()

to_delete = []
for m in measurements:
    if m['name'].startswith('xbos/'):
        to_delete.append(m)
        q = client.query('select * from "{0}"'.format(m['name']))
        col = m['name']
        count = 0
        for p in q.get_points():
            newp = {
                'tags': {
                    'collection': col,
                    'unit': p['unit'],
                    'name': p['name'],
                    'uuid': p['uuid'],
                    'prediction_step': p.get('prediction_step', None),
                },
                'measurement': 'timeseries',
                'time': p['time'],
                'fields': {
                    'prediction_time': p.get('prediction_time', None),
                    'value': float(p['value'])
                }
            }
            client.write_points([newp])
            count += 1
        print("Wrote {0} points from {1}".format(count, col))

print("\nCheck the 'timeseries' collection and then remove the following")                
for measurement in to_delete:
    print(measurement['name'])

