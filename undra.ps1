#!/usr/bin/env pwsh
$cwd = ("$PWD").replace("\", "/")

chdir "$cwd"
& "$PSScriptRoot\undra_start.exe" $args