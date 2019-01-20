# escape=`
FROM plugins/base:windows-amd64

LABEL maintainer="Drone.IO Community <drone-dev@googlegroups.com>" `
  org.label-schema.name="Drone Slack" `
  org.label-schema.vendor="Drone.IO Community" `
  org.label-schema.schema-version="1.0"

ADD release\drone-slack.exe c:\drone-slack.exe
ENTRYPOINT [ "c:\\drone-slack.exe" ]
