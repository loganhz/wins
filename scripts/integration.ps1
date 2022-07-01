#Requires -Version 5.0

$ErrorActionPreference = 'Stop'

Import-Module -WarningAction Ignore -Name "$PSScriptRoot\utils.psm1"

$SRC_PATH = (Resolve-Path "$PSScriptRoot\..").Path
Push-Location $SRC_PATH

Invoke-Expression -Command "$SRC_PATH\scripts\build.ps1"
if ($LASTEXITCODE -ne 0) {
    Log-Fatal "build failed prior to starting integration tests"
    exit $LASTEXITCODE
}
Invoke-Expression -Command "$SRC_PATH\tests\integration\integration_suite_test.ps1"
if ($LASTEXITCODE -ne 0) {
    Log-Fatal "integration test failed"
    exit $LASTEXITCODE
}
Write-Host -ForegroundColor Green "integration.ps1 has completed successfully."

Pop-Location
