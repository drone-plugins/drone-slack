FROM microsoft/nanoserver:1803
USER ContainerAdministrator
ADD release/windows/amd64/drone-slack.exe /drone-slack.exe
SHELL ["powershell", "-Command", "$ErrorActionPreference = 'Stop'; $ProgressPreference = 'SilentlyContinue';"]
ENTRYPOINT [ "\\drone-slack.exe" ]
