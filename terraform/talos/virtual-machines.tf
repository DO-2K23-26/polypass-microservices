resource "proxmox_virtual_environment_vm" "this" {
  for_each = var.nodes

  node_name = each.value.host_node

  name        = each.key
  description = each.value.machine_type == "controlplane" ? "Talos Control Plane" : "Talos Worker"
  tags        = each.value.machine_type == "controlplane" ? ["k8s", "talos", "polypass", "control-plane"] : ["k8s", "talos", "polypass", "worker"]
  on_boot     = true
  vm_id       = each.value.vm_id

  machine       = "q35"
  scsi_hardware = "virtio-scsi-single"
  bios          = "seabios"

  agent {
    enabled = true
  }

  cpu {
    cores = each.value.cpu
    type  = "host"
  }

  memory {
    dedicated = each.value.ram_dedicated
  }

  network_device {
    bridge      = "vmbr50"
    mac_address = each.value.mac_address
  }

  dynamic "network_device" {
    for_each = each.value.machine_type == "controlplane" ? [1] : []
    content {
      bridge      = "vmbr0"
      mac_address = each.value.secondary_mac_address
    }
  }

  disk {
    datastore_id = each.value.datastore_id
    interface    = "scsi0"
    iothread     = true
    cache        = "writethrough"
    discard      = "on"
    ssd          = true
    file_format  = "raw"
    size         = 20
    file_id      = proxmox_virtual_environment_download_file.this["${each.value.host_node}_${each.value.update == true ? local.update_image_id : local.image_id}"].id
  }

  boot_order = ["scsi0"]

  operating_system {
    type = "l26" # Linux Kernel 2.6 - 6.X.
  }

  dynamic "initialization" {
    for_each = each.value.machine_type == "controlplane" ? [1] : []
    content {
      datastore_id = each.value.datastore_id
      ip_config {
        ipv4 {
          address = "${each.value.ip}/24"
          gateway = var.cluster.private_gateway
        }
      }
      ip_config {
        ipv4 {
          address = "${var.cluster.endpoint}/24"
          gateway = var.cluster.gateway
        }
      }
    }
  }

  dynamic "initialization" {
    for_each = each.value.machine_type == "worker" ? [1] : []
    content {
      datastore_id = each.value.datastore_id
      ip_config {
        ipv4 {
          address = "${each.value.ip}/24"
          gateway = var.cluster.private_gateway
        }
      }
    }
  }

  dynamic "hostpci" {
    for_each = each.value.igpu ? [1] : []
    content {
      # Passthrough iGPU
      device  = "hostpci0"
      mapping = "iGPU"
      pcie    = true
      rombar  = true
      xvga    = false
    }
  }
}
