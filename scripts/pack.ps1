#Requires -Version 5.1
<#
.SYNOPSIS
    HyperSKU CLI 打包脚本

.DESCRIPTION
    1. 编译 hypersku-cli 二进制文件，并打包到 ~\.hypersku-cli
    2. 将 skills/ 下每个能力域目录安装为独立 skill 到 ~\.agents\skills\（SKILL.md + references/）

.EXAMPLE
    powershell -ExecutionPolicy Bypass -File scripts\pack.ps1
#>

$ErrorActionPreference = 'Stop'

# 让控制台正确显示中文
try { [Console]::OutputEncoding = [System.Text.Encoding]::UTF8 } catch { }

# 项目根目录（脚本所在目录的上一级）
$ProjectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $ProjectRoot

# ========== 配置 ==========
$AppName = 'hypersku-cli'
$Version = '0.1.0'

# 读取 git 提交哈希（非 git 环境时回退为 unknown）
$Commit = 'unknown'
if (Get-Command git -ErrorAction SilentlyContinue) {
    $Commit = & git rev-parse --short HEAD 2>$null
    if (-not $Commit) { $Commit = 'unknown' }
}
$Date = (Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')

# 目标目录
$BinDir = Join-Path $HOME '.hypersku-cli/bin'

Write-Host ''
Write-Host '==============================' -ForegroundColor Cyan
Write-Host '  HyperSKU CLI 打包' -ForegroundColor Cyan
Write-Host '==============================' -ForegroundColor Cyan

# ========== 1. 编译二进制 ==========
Write-Host ''
Write-Host '[1/2] 编译二进制文件 ...' -ForegroundColor Yellow
$LdFlags = "-X github.com/hypersku/hypersku-cli/internal/version.Version=$Version" +
           " -X github.com/hypersku/hypersku-cli/internal/version.Commit=$Commit" +
           " -X github.com/hypersku/hypersku-cli/internal/version.Date=$Date"
$BuildDir = Join-Path $ProjectRoot 'build'
New-Item -ItemType Directory -Force -Path $BuildDir | Out-Null
& go build "-ldflags=$LdFlags" -o (Join-Path $BuildDir "$AppName.exe") .
if ($LASTEXITCODE -ne 0) { throw 'go build 编译失败' }

# 复制二进制到 ~\.hypersku-cli
Write-Host '[1/2] 复制二进制到 ~\.hypersku-cli ...' -ForegroundColor Yellow
New-Item -ItemType Directory -Force -Path $BinDir | Out-Null
$BinFile = Join-Path $BinDir "$AppName.exe"
Copy-Item -Path (Join-Path $BuildDir "$AppName.exe") -Destination $BinFile -Force
Write-Host ("  OK 二进制已安装: " + $BinFile) -ForegroundColor Green

# ========== 2. 打包 skills ==========
Write-Host ''
Write-Host '[2/2] 安装 skills 到 ~\.agents\skills\（每个能力域一个独立 skill） ...' -ForegroundColor Yellow
$SourceSkills = Join-Path $ProjectRoot 'skills'
$AgentsSkillsRoot = Join-Path $HOME '.agents\skills'

New-Item -ItemType Directory -Force -Path $AgentsSkillsRoot | Out-Null

# 清理本工具旧安装痕迹：旧版统一的 ehub 目录 + 旧无前缀能力域目录 + 当前前缀目录（避免残留旧结构）
$OldDirs = @('ehub',
             'after-sales', 'after-sales-apply', 'customer', 'customer-profile-analysis',
             'domestic-exception-handling', 'domestic-third-trade-exception', 'logistics', 'purchase', 'warehouse',
             'hypersku-auth', 'hypersku-cli',
             'hypersku-after-sales', 'hypersku-after-sales-apply', 'hypersku-customer', 'hypersku-customer-profile-analysis',
             'hypersku-domestic-exception-handling', 'hypersku-domestic-third-trade-exception', 'hypersku-logistics',
             'hypersku-purchase', 'hypersku-warehouse')
foreach ($d in $OldDirs) {
    $target = Join-Path $AgentsSkillsRoot $d
    if (Test-Path $target) { Remove-Item -Path $target -Recurse -Force }
}

# 复制 skills 下每个能力域目录为独立 skill（SKILL.md + references/）
$InstalledSkills = @()
Get-ChildItem -Path $SourceSkills -Directory | ForEach-Object {
    if (-not (Test-Path (Join-Path $_.FullName 'SKILL.md'))) { return }
    $dest = Join-Path $AgentsSkillsRoot $_.Name
    Copy-Item -Path $_.FullName -Destination $dest -Recurse -Force
    $InstalledSkills += $_.Name
}
Write-Host ("  OK skills 已安装: " + $AgentsSkillsRoot) -ForegroundColor Green
$InstalledSkills | ForEach-Object { Write-Host ("  - " + $_) }

# ========== 完成 ==========
Write-Host ''
Write-Host '打包完成 OK' -ForegroundColor Green
Write-Host ("  - 二进制: " + $BinFile)
Write-Host ("  - Skills: " + $AgentsSkillsRoot)
