# Woodpecker CI Slack Plugin

Woodpecker plugin for sending Slack notifications cloned from: https://github.com/drone-plugins/drone-slack

For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-slack/).

## Build

Build the binary with the following commands:

```
mise build
```

This will produce a file called `drone-slack` in your root directory.

## Docker

Build the Docker image and test-run it:

```
mise docker:run
```

This will produce docker image `gowerstreet/slack:local` and run it with a few test variables. You should be able to see the output in the `#alerts` slack channel.
