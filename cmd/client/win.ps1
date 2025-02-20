# $OS = "darwin"
# $ARCH = "amd64"
# $S = $OS +":"+$ARCH
# Write-Output $S
# $Env:GOOS = $OS
# $Env:GOARCH = $ARCH
# go build -ldflags "-X main.Version=v1.0 -X 'main.BuildTime=$(Get-Date)'"

# $OS = "linux"
# $ARCH = "amd64"
# $S = $OS +":"+ $ARCH
# Write-Output $S
# $Env:GOOS = $OS
# $Env:GOARCH = $ARCH
# go build -ldflags "-X main.Version=v1.0 -X 'main.BuildTime=$(Get-Date)'"

$OS = "windows"
# $ARCH = "386"
$S = $OS +":" + $ARCH
Write-Output $S
$Env:GOOS = $OS
$Env:GOARCH = $ARCH
go build -ldflags "-X main.Version=v1.0 -X 'main.BuildTime=$(Get-Date)'"

.\client.exe -config="local_cfg.yaml"
