#!/usr/bin/env pwsh
$cwd = ("$PWD").replace("\", "/")
& "$PSScriptRoot\undra_start.exe" $args -cwd="$cwd"