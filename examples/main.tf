
resource "launchpad_rocket" "example1" {
    name = "example1"
}

resource "launchpad_rocket" "example2" {
    name = "example2"
}

resource "launchpad_mission" "example1" {
    name = "example1"
    rockets = [
        {
            id = launchpad_rocket.example1.id
            config = {}
        },
        {
            id = launchpad_rocket.example2.id
            config = {}
        }
    ]
}

resource "launchpad_mission" "example2" {
    name = "example2"
    rockets = [
        {
            id = launchpad_rocket.example2.id
            config = {}
        }
    ]
}
