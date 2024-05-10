---
page_title: "Scaleway Zones and Regions"
---

# Scaleway Zones and Regions

Scaleway's products are deployed across multiple datacenter in the world.

For technical and legal reasons, some products are splitted by Region or by Availability Zones.
When using such product, you can choose the location that better fits your need (country, latency, ...).

## Regions

A Region is represented as a Geographical area such as France (Paris: `fr-par`) or the Netherlands (Amsterdam: `nl-ams`).
It can contain multiple Availability Zones.


## Zones

In order to deploy highly available application, a region can be divided in many Availability Zones (AZ).
Latency between multiple AZ of the same region are low as they have a common network layer.

List of availability zones by regions:

- France - Paris (`fr-par`)
    - `fr-par-1`
    - `fr-par-2`
    - `fr-par-3`
- The Netherlands - Amsterdam (`nl-ams`)
    - `nl-ams-1`
    - `nl-ams-2`
    - `nl-ams-3`
- Poland - Warsaw (`pl-waw`)
    - `pl-waw-1`
    - `pl-waw-2`
    - `pl-waw-3`

## Resource IDs

To save this notion of regions and zones in the state, all the Terraform IDs of Scaleway contain the region or zone.
This is saved in the following format: `{zone|region}/{resource_id}`.
Where `zone` or `region` is the place where the resource is created and where `resource_id` is the ID that is used on Scaleway's console/API.

If you need to retrieve the raw ID of the resource, you can either :

- use the `trimprefix` function :

`id = trimprefix(scaleway_resource.name.id, "${scaleway_resource.name.zone|region}/")`

- use the `split` function :

`zone|region = split("/", scaleway_resource.name.id)[0]`

`id = split("/", scaleway_resource.name.id)[1]`

---

More information regarding zones and regions can be found [here](https://www.scaleway.com/en/developers/api/#regions-and-zones).
