plugin "terraform" {
  enabled = false
  version = "0.2.1"
  source  = "github.com/terraform-linters/tflint-ruleset-terraform"
}

config {
  module = true
}
