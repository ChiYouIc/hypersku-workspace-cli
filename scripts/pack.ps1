#Requires -Version 5.1
<#
.SYNOPSIS
    HyperSKU CLI 打包脚本

.DESCRIPTION
    1. 编译 hypersku-cli 二进制文件，并打包到 ~\.hypersku-cli
    2. 将 skills 目录打包到 ~\.agents\skills\ehub（重命名为 ehub）

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
$BinDir    = Join-Path $HOME '.hypersku-cli'
$SkillsDir = Join-Path $HOME '.agents\skills\ehub'

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
Write-Host '[2/2] 打包 skills 到 ~\.agents\skills\ehub ...' -ForegroundColor Yellow
$SourceSkills = Join-Path $ProjectRoot 'skills'

# 先清空旧目标，避免残留旧文件
if (Test-Path $SkillsDir) {
    Get-ChildItem -Path $SkillsDir -Force | Remove-Item -Recurse -Force
}
New-Item -ItemType Directory -Force -Path $SkillsDir | Out-Null

# 复制 skills 下所有内容（SKILL.md 及各个子 skill 目录）
Copy-Item -Path (Join-Path $SourceSkills '*') -Destination $SkillsDir -Recurse -Force
Write-Host ("  OK skills 已安装: " + $SkillsDir) -ForegroundColor Green

# ========== 完成 ==========
Write-Host ''
Write-Host '打包完成 OK' -ForegroundColor Green
Write-Host ("  - 二进制: " + $BinFile)
Write-Host ("  - Skills: " + $SkillsDir)
