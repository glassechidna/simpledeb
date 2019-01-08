workflow "goreleaser" {
  on = "push"
  resolves = ["docker://goreleaser/goreleaser"]
}

action "Only publish on tag" {
  uses = "actions/bin/filter@b2bea0749eed6beb495a8fa194c071847af60ea1"
  args = "tag v*"
}

action "docker://goreleaser/goreleaser" {
  uses = "docker://goreleaser/goreleaser"
  needs = ["Only publish on tag"]
}
