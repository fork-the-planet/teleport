resource "teleport_autoupdate_config" "test" {
  version = "v1"
  spec = {
    tools = {
      mode = "enabled"
    }
    agents = {
      mode     = "enabled"
      strategy = "halt-on-error"
      schedules = {
        regular = [
          # dev is updated at 4:00 UTC
          { name = "dev", days = ["Mon", "Tue", "Wed", "Thu"], start_hour : 4 },
          # staging is updated at 14:00 UTC
          { name = "staging", days = ["Mon", "Tue", "Wed", "Thu"], start_hour : 14 },
          # prod is updated at 14:00 UTC the next day
          { name = "prod", days = ["Mon", "Tue", "Wed", "Thu"], start_hour : 14, wait_hours : 24 },
        ]
      }
    }
  }
}
