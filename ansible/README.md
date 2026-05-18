# Ansible Automation

This folder demonstrates the configuration management and automated deployment requirement.

## Configure

Copy the inventory and set your VM IP/user:

```bash
cp ansible/inventory.example.ini ansible/inventory.ini
```

Edit `ansible/inventory.ini`, then run the playbook with your real repository URL:

```bash
ansible-playbook -i ansible/inventory.ini ansible/playbook.yml \
  -e repo_url=https://github.com/YOUR_ACCOUNT/YOUR_REPO.git \
  -e repo_version=main
```

## What It Does

- Installs Docker Engine and Docker Compose plugin.
- Clones or updates the project repository.
- Writes the required `.env` file.
- Validates `docker compose` syntax.
- Runs `docker compose up -d --build`.
- Prints `docker compose ps` output for evidence.

Screenshot evidence:

- Successful Ansible play recap.
- `TASK [Print deployment status]` output.
- Application open through the VM public IP.
