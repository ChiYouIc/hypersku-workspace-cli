#Requires -Version 5.1
<#
.SYNOPSIS
    HyperSKU CLI workbuddy 平台打包脚本

.DESCRIPTION
    专门用于打包 workbuddy 平台可识别的发布物。同一份源（skills/ + connector/ + Go 源码），
    按 -Target 参数输出两类产物：

      workbuddy-skill      每个能力域一个独立 skill zip → build\dist\workbuddy-skill\
                           （zip 根 = SKILL.md + references/，符合 workbuddy 两级结构规范）
      workbuddy-connector  完整连接器 zip（connector-meta.json + cli.json + icon.svg +
                           skills/ 全部子目录 + 三平台二进制）→ build\dist\workbuddy-connector\
      all                  依次执行以上全部目标（默认）

    本地开发安装请使用 scripts/pack.ps1（make pack），不属于 workbuddy 发布物，故不在此脚本内。

    未来兼容其它 AI 工作台：新建独立脚本，源数据不变。

.PARAMETER Target
    打包目标：workbuddy-skill | workbuddy-connector | all（默认 all）

.PARAMETER Version
    版本号（写入二进制与连接器元数据；默认取 Makefile 同步值 0.1.0）

.PARAMETER SkipBinary
    连接器包不附带平台二进制（纯 skills + 元数据预览场景）

.EXAMPLE
    powershell -ExecutionPolicy Bypass -File scripts\pack-workbuddy.ps1 -Target workbuddy-skill

.EXAMPLE
    powershell -ExecutionPolicy Bypass -File scripts\pack-workbuddy.ps1 -Target workbuddy-connector

.EXAMPLE
    powershell -ExecutionPolicy Bypass -File scripts\pack-workbuddy.ps1 -Target all
#>

param(
    [ValidateSet('workbuddy-skill', 'workbuddy-connector', 'all')]
    [string]$Target = 'all',
    [string]$Version = '0.1.0',
    [switch]$SkipBinary
)

$ErrorActionPreference = 'Stop'

# 让控制台正确显示中文
try { [Console]::OutputEncoding = [System.Text.Encoding]::UTF8 } catch { }

# 项目根目录（脚本所在目录的上一级）
$ProjectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $ProjectRoot

# ========== 公共配置 ==========
$AppName         = 'hypersku-cli'
$SkillsSource    = Join-Path $ProjectRoot 'skills'
$ConnectorSource = Join-Path $ProjectRoot 'connector'
$DistRoot        = Join-Path $ProjectRoot 'build\dist'

# 读取 git 提交哈希（非 git 环境时回退为 unknown）
$Commit = 'unknown'
if (Get-Command git -ErrorAction SilentlyContinue) {
    $Commit = & git rev-parse --short HEAD 2>$null
    if (-not $Commit) { $Commit = 'unknown' }
}
$Date = (Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')
$LdFlags = "-X github.com/hypersku/hypersku-cli/internal/version.Version=$Version" +
           " -X github.com/hypersku/hypersku-cli/internal/version.Commit=$Commit" +
           " -X github.com/hypersku/hypersku-cli/internal/version.Date=$Date"

# ========== 工具函数 ==========

function Get-SkillDirs {
    # 返回 skills/ 下所有含 SKILL.md 的能力域目录
    Get-ChildItem -Path $SkillsSource -Directory | Where-Object {
        Test-Path (Join-Path $_.FullName 'SKILL.md')
    }
}

function Build-Binaries {
    # 编译三平台二进制到 build\bin\
    param([string]$OutDir)

    $platforms = @(
        @{ GOOS = 'windows'; GOARCH = 'amd64'; Name = "$AppName-win32-amd64.exe" },
        @{ GOOS = 'linux';   GOARCH = 'amd64'; Name = "$AppName-linux-amd64" },
        @{ GOOS = 'darwin';  GOARCH = 'amd64'; Name = "$AppName-darwin-amd64" }
    )
    New-Item -ItemType Directory -Force -Path $OutDir | Out-Null
    foreach ($p in $platforms) {
        Write-Host ("  编译 " + $p.Name + " ...") -ForegroundColor Yellow
        $env:GOOS = $p.GOOS
        $env:GOARCH = $p.GOARCH
        & go build "-ldflags=$LdFlags" -o (Join-Path $OutDir $p.Name) .
        if ($LASTEXITCODE -ne 0) { throw ("go build " + $p.Name + " 失败") }
    }
    Remove-Item Env:GOOS, Env:GOARCH -ErrorAction SilentlyContinue
}

function Reset-Dir {
    param([string]$Path)
    if (Test-Path $Path) { Remove-Item -Path $Path -Recurse -Force }
    New-Item -ItemType Directory -Force -Path $Path | Out-Null
}

# ========== 目标：workbuddy-skill（每能力域一个独立 zip） ==========

function Build-WorkbuddySkill {
    $outDir = Join-Path $DistRoot 'workbuddy-skill'
    Reset-Dir $outDir

    Write-Host ''
    Write-Host '>>> [workbuddy-skill] 打包独立技能包（每能力域一个 zip）' -ForegroundColor Cyan

    $packed = @()
    foreach ($dir in (Get-SkillDirs)) {
        $zipPath = Join-Path $outDir ($dir.Name + '.zip')
        # zip 根 = SKILL.md + references/（两级结构，符合平台规范）
        Compress-Archive -Path (Join-Path $dir.FullName '*') -DestinationPath $zipPath -CompressionLevel Optimal -Force
        Write-Host ("  OK " + $dir.Name + ".zip") -ForegroundColor Green
        $packed += $dir.Name
    }
    Write-Host ("  共 " + $packed.Count + " 个独立 skill 包 → " + $outDir) -ForegroundColor Green
}

# ========== 目标：workbuddy-connector（完整连接器 zip） ==========

function Build-WorkbuddyConnector {
    $outDir = Join-Path $DistRoot 'workbuddy-connector'
    Reset-Dir $outDir

    Write-Host ''
    Write-Host '>>> [workbuddy-connector] 打包完整连接器（元数据 + skills + 可选二进制）' -ForegroundColor Cyan

    # 校验连接器必备文件
    foreach ($f in @('connector-meta.json', 'cli.json', 'icon.svg')) {
        if (-not (Test-Path (Join-Path $ConnectorSource $f))) {
            throw ("connector/" + $f + " 缺失，无法打包连接器")
        }
    }

    # 暂存目录组装连接器结构：根 = 三件套 + skills/
    $stage = Join-Path $outDir 'stage'
    New-Item -ItemType Directory -Force -Path $stage | Out-Null

    Copy-Item -Path (Join-Path $ConnectorSource '*') -Destination $stage -Recurse -Force
    $stageSkills = Join-Path $stage 'skills'
    New-Item -ItemType Directory -Force -Path $stageSkills | Out-Null
    foreach ($dir in (Get-SkillDirs)) {
        Copy-Item -Path $dir.FullName -Destination (Join-Path $stageSkills $dir.Name) -Recurse -Force
    }

    # 同步版本号到 connector-meta.json（保持发布物版本与二进制一致）
    # 注意：PS 5.1 的 Get-Content 默认 ANSI，必须显式指定 UTF-8 读取
    $metaPath = Join-Path $stage 'connector-meta.json'
    $meta = [IO.File]::ReadAllText($metaPath, [Text.Encoding]::UTF8) | ConvertFrom-Json
    $meta.version = $Version
    [IO.File]::WriteAllText($metaPath, ($meta | ConvertTo-Json -Depth 10), (New-Object Text.UTF8Encoding $false))

    # 可选：附带三平台二进制（离线分发场景；线上分发走 init 下载，见 cli.json）
    if (-not $SkipBinary) {
        $binDir = Join-Path $ProjectRoot 'build\bin'
        Build-Binaries -OutDir $binDir
        $stageBin = Join-Path $stage 'bin'
        New-Item -ItemType Directory -Force -Path $stageBin | Out-Null
        Get-ChildItem -Path $binDir -File | ForEach-Object {
            Copy-Item -Path $_.FullName -Destination (Join-Path $stageBin $_.Name) -Force
        }
    }

    $zipPath = Join-Path $outDir 'hypersku-cli-connector.zip'
    Compress-Archive -Path (Join-Path $stage '*') -DestinationPath $zipPath -CompressionLevel Optimal -Force
    Remove-Item -Path $stage -Recurse -Force
    Write-Host ("  OK hypersku-cli-connector.zip → " + $outDir) -ForegroundColor Green
}

# ========== 目标注册表（新增打包目标在此登记） ==========
$Targets = @{
    'workbuddy-skill'     = ${function:Build-WorkbuddySkill}
    'workbuddy-connector' = ${function:Build-WorkbuddyConnector}
}

# ========== 执行 ==========
Write-Host ''
Write-Host '==============================' -ForegroundColor Cyan
Write-Host ("  HyperSKU CLI workbuddy 打包  target=" + $Target) -ForegroundColor Cyan
Write-Host '==============================' -ForegroundColor Cyan

$toRun = if ($Target -eq 'all') { @('workbuddy-skill', 'workbuddy-connector') } else { @($Target) }
foreach ($t in $toRun) {
    & $Targets[$t]
}

Write-Host ''
Write-Host 'workbuddy 打包完成 OK' -ForegroundColor Green
