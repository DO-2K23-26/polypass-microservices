module "talos" {
  source = "./talos"

  providers = {
    proxmox = proxmox
  }

  image = {
    version = "v1.10.0"
    schematic = file("${path.module}/talos/image/schematic.yaml")
  }

  cluster = {
    name            = "polypass"
    endpoint        = "162.38.112.168"
    gateway         = "162.38.112.254"
    private_gateway = "10.15.50.1"
    talos_version   = "v1.10.0"
    proxmox_cluster = "serpentard"
  }

  nodes = {
    "polypass-ctrl-00" = {
      host_node     = "serpentard-1"
      machine_type  = "controlplane"
      ip            = "10.15.50.110"
      mac_address   = "BC:24:11:2E:C8:00"
      secondary_mac_address = "BC:24:11:2E:C8:01"
      vm_id         = 850
      cpu           = 4
      ram_dedicated = 8192
    }
    "polypass-work-00" = {
      host_node     = "serpentard-2"
      machine_type  = "worker"
      ip            = "10.15.50.111"
      mac_address   = "BC:24:11:2E:C8:02"
      vm_id         = 851
      cpu           = 8
      ram_dedicated = 16384
    }
    "polypass-work-01" = {
      host_node     = "serpentard-3"
      machine_type  = "worker"
      ip            = "10.15.50.112"
      mac_address   = "BC:24:11:2E:C8:03"
      vm_id         = 852
      cpu           = 8
      ram_dedicated = 16384
    }
    "polypass-work-02" = {
      host_node     = "serpentard-1"
      machine_type  = "worker"
      ip            = "10.15.50.113"
      mac_address   = "BC:24:11:2E:C8:04"
      vm_id         = 853
      cpu           = 8
      ram_dedicated = 16384
    }

  }
}
