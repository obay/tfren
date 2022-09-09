# Terraform Rename (tfren)
Terraform Rename (tfren) is a tool to rename Terraform files using the resource type and name.

Currently it runs only in the current directory and does not recurse into subdirectories.

`tfren` assumes you have one resource per file and uses the the first line (e.g. `resource "azurerm_virtual_network" "sandbox-vnet" {`) to build the file name (which will be `resource.azurerm_virtual_network.sandbox-vnet.tf` in this case).


## Installation
### Homebrew
```bash
brew install obay/tap/tfren
```

## Usage
Simply switch to the directory containing your Terraform files and run `tfren`.
```bash
tfren
```
