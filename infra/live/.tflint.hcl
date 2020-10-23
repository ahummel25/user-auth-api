plugin "terraform" {
  enabled = false
  version = "0.2.2"
  source  = "github.com/terraform-linters/tflint-ruleset-terraform"
}

config {
  module = true
}
