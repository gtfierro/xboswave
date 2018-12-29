# XBOS(2) 

- [Issues](https://todo.sr.ht/%7Egabe/xbos)
- [Docs](https://man.sr.ht/%7Egabe/xbos/)

## Ingester

How do get started:


1. Install [BtrDB dev machine](https://docs.smartgrid.store/development-environment.html) OR [InfluxDB](https://docs.influxdata.com/influxdb/v1.7/introduction/)
2. Install [WAVEMQ](https://github.com/immesys/wavemq):
    - can use ansible:
    1. [Install ansible](https://docs.ansible.com/ansible/2.7/installation_guide/intro_installation.html#installing-the-control-machine)
    2. Clone this repo
    3. Go to `ansible/`
    4. Run the playbook:
        ```
        ansible-playbook -K wavemq-playbook.yml
        ```
3. Create entity (docs coming)
    - read [wave docs](https://github.com/immesys/wave) for now
4. Give perms
    - read [wave docs](https://github.com/immesys/wave) for now
5. Build and run ingester:
    ```
    cd ingester
    make run
    ```
