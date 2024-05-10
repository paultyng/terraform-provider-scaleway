---
page_title: "Using Elastic Metal servers in a Kubernetes cluster"
---

# How to use Elastic Metal servers with a Kosmos multicloud cluster

In this guide you will learn how to deploy bare metal nodes on your Kubernetes cluster instead of Instances. To be able
to do this, you need to have a [Kosmos multicloud cluster](../resources/k8s_cluster.md#multicloud) instead of a Kapsule
cluster, and add your [pools](../resources/k8s_pool.md) as the "external" node type.
Once you have set up your infrastructure, you will have to run a program on your server so it is configured and registered
as a node by the apiserver. This can be achieved manually ([method A](#method-a-manually-via-ssh-connexion)), or you can automate the process ([method B](#method-b-fully-automated-with-terraform-remote-exec)).

## Prerequisites

* An SSH key, required to connect to your Elastic Metal server.
    * If you need help generating an SSH key, visit [this tutorial](https://www.scaleway.com/en/docs/identity-and-access-management/organizations-and-projects/how-to/create-ssh-key/).
    * You can view your SSH keys in the [console](https://console.scaleway.com/project/credentials) and add a new one.

## Setup

```hcl
###############################################
#         CONFIGURE THE K8S CLUSTER           #
###############################################

resource "scaleway_k8s_cluster" "multicloud" {
  name    = "multicloud-cluster"
  type    = "multicloud"
  version = "1.27.1"
  cni     = "kilo"
  region  = "fr-par"
  delete_additional_resources = false
}

resource "scaleway_k8s_pool" "pool" {
  cluster_id  = scaleway_k8s_cluster.multicloud.id
  name        = "multicloud-pool"
  node_type   = "external"
  size        = 0
  region      = "fr-par"
}

###############################################
#     CONFIGURE THE ELASTIC METAL SERVER      #
###############################################

# Select at least one SSH key to connect to your server
resource "scaleway_iam_ssh_key" "key" {
  name = "ssh-key"
  public_key = file("~/.ssh/id_ed25519.pub")
}
# Select the type of offer for your server
data "scaleway_baremetal_offer" "offer" {
  name = "EM-B112X-SSD"
}
# Select the OS you want installed on your server
data "scaleway_baremetal_os" "os" {
  name = "Ubuntu"
  version = "20.04 LTS (Focal Fossa)"
}

resource "scaleway_baremetal_server" "server" {
  offer       = data.scaleway_baremetal_offer.offer.name  # The name of the Elastic Metal offer
  os          = data.scaleway_baremetal_os.os.id          # The ID of the OS
  ssh_key_ids = [scaleway_iam_ssh_key.key.id]             # The list of SSH key IDs allowed to connect to the server
  zone        = "fr-par-2"
}
```

### Notes

* There is also an ARM binary (named `node-agent_linux_arm64`) for ARM based nodes.
* If you want a fully automated process, don't apply this configuration yet because you will have to modify the spec of
the server, which will trigger a new installation process that will take some time. You should instead apply this [configuration](#method-b-fully-automated-with-terraform-remote-exec)
* `kilo` is the only CNI compatible with multicloud clusters
* In this example, we use data sources to fill the `offer` and `os` fields of the bare metal server rather than giving
raw values because it allows to check their availability in the zone you want to work with before provisioning the server
    * See the [baremetal offer datasource](../data-sources/baremetal_offer.md)
    * See the [baremetal os datasource](../data-sources/baremetal_os.md)
    * For more information on bare metal servers specs, visit the [resource documentation](../resources/baremetal_server.md)
* If you want to link any already existing resource, you can import it to the Terraform state by running :

   ```bash
   terraform import <scaleway_resource_type> <locality>/<id>
   ```

## Configure the Elastic Metal server

### Method A: Manually via SSH connexion

1. Get your server's public IP and SSH to the server :

    ```bash
    ssh <user>@<baremetal_server_ip>
    ```

2. Download the node-agent program :

    ```bash
    wget https://scwcontainermulticloud.s3.fr-par.scw.cloud/node-agent_linux_amd64 && chmod +x node-agent_linux_amd64
    ```

3. Export the required environment variables :

    ```bash
    export POOL_ID=<pool_id>  POOL_REGION=<cluster_region>  SCW_SECRET_KEY=<secret_key>
    ```

4. Execute the program to attach the node to the multicloud pool :

    ```bash
    sudo ./node-agent_linux_amd64 -loglevel 0 -no-controller
    ```

### Method B: Fully automated with Terraform "remote-exec"

If you want the configuration process to be automated, you will have to add a few things to your server configuration.
In addition to the SSH key file and the data sources for the offer and the os, add your secret key and give the
configuration instructions in the bare metal server spec.

```hcl
# Put your secret key in a file on your local machine
data "local_sensitive_file" "secret_key" {
    filename = pathexpand("~/path/to/secret/key")
}

resource "scaleway_baremetal_server" "server" {
    offer       = data.scaleway_baremetal_offer.offer.name
    os          = data.scaleway_baremetal_os.os.id
    ssh_key_ids = [scaleway_iam_ssh_key.key.id]

    # Configure the SSH connexion used by Terraform for the remote execution  
    connection {
      type     = "ssh"
      user     = "ubuntu"
      host     = one([for k in self.ips : k if k.version == "IPv4"]).address   # We look for the IPv4 in the list of IPs
    }

    # Download and execute the configuration script
    provisioner "remote-exec" {
      inline = [
        "wget https://scwcontainermulticloud.s3.fr-par.scw.cloud/node-agent_linux_amd64 > log && chmod +x node-agent_linux_amd64",
        "echo \"\nPOOL_ID=${split("/", scaleway_k8s_pool.pool.id)[1]}\nPOOL_REGION=${scaleway_k8s_pool.pool.region}\nSCW_SECRET_KEY=${data.local_sensitive_file.secret_key.content}\" >> log",
        "export POOL_ID=${split("/", scaleway_k8s_pool.pool.id)[1]}  POOL_REGION=${scaleway_k8s_pool.pool.region}  SCW_SECRET_KEY=${data.local_sensitive_file.secret_key.content}",
        "sudo ./node-agent_linux_amd64 -loglevel 0 -no-controller >> log",
      ]
    }
}
```

### Verify the installation

You check that everything went well by :

* In the CLI : by listing the nodes of the cluster|pool with the CLI

    ```
    scw k8s node list cluster-id=<cluster_id> [pool-id=<pool_id>]
    ```

* In the console : by checking the `Nodes` tab of your cluster.
* On the server : by connecting via SSH and checking the `log` file located in `/home/ubuntu`. It should display the
following lines :

```
[...]
{"time":"2023-05-24T16:47:38.750041045Z","level":"DEBUG","msg":"writing kubelet config and CA"}
{"time":"2023-05-24T16:47:38.750272845Z","level":"DEBUG","msg":"writing kubelet env file and systemd service"}
{"time":"2023-05-24T16:47:38.750410725Z","level":"DEBUG","msg":"writing kubeconfig files"}
{"time":"2023-05-24T16:47:38.751171166Z","level":"DEBUG","msg":"starting containerd systemd service"}
{"time":"2023-05-24T16:47:39.042781246Z","level":"DEBUG","msg":"starting kubelet systemd service"}
{"time":"2023-05-24T16:47:39.056392423Z","level":"INFO","msg":"successfully started kubelet"}
```

  If something went wrong you should be able to find useful information for troubleshooting in here, like the
  environment values that got exported. You can also rerun the command with a different loglevel (lower is more verbose).
