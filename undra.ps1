#!/usr/bin/env pwsh
$cwd = ("$PWD").replace("\", "/")

chdir "$cwd"
& "$PSScriptRoot\undrastart.exe" $args