---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "ceph_auth Resource - terraform-provider-ceph"
subcategory: ""
description: |-
  
---

# ceph_auth (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `entity` (String) The entity name (i.e.: client.admin)

### Optional

- `caps` (String) The caps wanted for the entity

### Read-Only

- `id` (String) The ID of this resource.
- `key` (String) The cephx key of the entity
- `keyring` (String) The cephx keyring of the entity

