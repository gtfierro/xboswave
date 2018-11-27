# XBOS(2) 

How do get started:


### Ansible

1. [Install ansible](https://docs.ansible.com/ansible/2.7/installation_guide/intro_installation.html#installing-the-control-machine)
2. Clone this repo
3. Go to `ansible/`
4. Run the playbook:
    ```
    ansible-playbook -K wavemq-playbook.yml
    ```

### Do it yourself

1. Install [BtrDB dev machine](https://docs.smartgrid.store/development-environment.html)
2. Install [WAVEMQ](https://github.com/immesys/wavemq)
3. Create entity for 
4. Build and run ingester:
    ```
    cd ingester
    make run
    ```
