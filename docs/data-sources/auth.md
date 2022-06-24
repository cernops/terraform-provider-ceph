---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ceph_auth Data Source - terraform-provider-ceph"
subcategory: ""
description: |-
  This data source allows you to get information about a ceph client.
---

# ceph_auth (Data Source)

This data source allows you to get information about a ceph client.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `entity` (String) The entity name (i.e.: client.admin)

### Read-Only

- `caps` (Map of String) The caps of the entity
- `id` (String) The ID of this resource.
- `key` (String) The cephx key of the entity
- `keyring` (String) The cephx keyring of the entity

