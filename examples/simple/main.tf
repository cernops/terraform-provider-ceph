terraform {
  required_version = ">= 0.14.0"
  required_providers {
    ceph = {
      source  = "cernops/ceph"
      version = "~> 0.1.0"
    }
  }
}

provider "ceph" {
  entity = "client.admin"
}

resource "ceph_wait_online" "wait" {
    cluster_name = "my-super-cluster"
}

resource "ceph_auth" "test" {
    entity = "client.test"
    caps = {
        "mon": "allow *",
        "osd": "allow rw",
        "mds": "allow rw"
    }
    depends_on = [
        ceph_wait_online.test
    ]
}
