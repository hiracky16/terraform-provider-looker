---
page_title: "looker_content_metadata_access Resource - terraform-provider-looker"
subcategory: ""
description: |-
  Manages Looker Content Metadata Access (permissions for Folders, Dashboards, etc.).
---

# looker_content_metadata_access (Resource)

Manages Looker Content Metadata Access. This can be used to set permissions (view or edit) for folders, dashboards, and other content types for specific users or groups.

## Example Usage

```terraform
resource "looker_folder" "my_folder" {
  name      = "My Custom Folder"
  parent_id = "1"
}

resource "looker_group" "my_group" {
  name = "My Group"
}

resource "looker_content_metadata_access" "folder_view" {
  content_metadata_id = looker_folder.my_folder.content_metadata_id
  group_id            = looker_group.my_group.id
  permission_type     = "view"
}
```

## Schema

### Required

- `content_metadata_id` (String) The ID of the content metadata to which access is being granted.
- `permission_type` (String) The type of permission to grant. Valid values are `view` and `edit`.

### Optional

- `group_id` (String) The ID of the group to which access is being granted. Exactly one of `user_id` or `group_id` must be provided.
- `user_id` (String) The ID of the user to which access is being granted. Exactly one of `user_id` or `group_id` must be provided.

### Read-Only

- `id` (String) The ID of this resource.
