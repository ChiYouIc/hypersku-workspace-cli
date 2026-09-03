#Requires -Version 5.1
<#
.SYNOPSIS
    HyperSKU CLI 技能打包上传脚本（多包模式）

.DESCRIPTION
    将 skills/ 下每个能力域目录（SKILL.md + references/）分别压缩为独立 zip 包到 build/skills/ 目录。
    每个 zip 的根目录直接包含 SKILL.md 及 references/，不带外层能力域文件夹，符合平台两级目录限制。
    可选：上传到百炼 Agent Studio（每个 skill 一个 skill_id，见 $SkillIds 映射表）。

.EXAMPLE
    powershell -ExecutionPolicy Bypass -File scripts\upload_skill.ps1
#>

$ErrorActionPreference = 'Stop'

# 让控制台正确显示中文
try { [Console]::OutputEncoding = [System.Text.Encoding]::UTF8 } catch { }

# 项目根目录（脚本所在目录的上一级）
$ProjectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $ProjectRoot

# ========== 配置 ==========
$SkillsSource    = Join-Path $ProjectRoot 'skills'
$OutDir          = Join-Path $ProjectRoot 'build\skills'
$DashscopeApiUrl = $env:AGENTSTUDIO_URL
$DashscopeApiKey = $env:BAILIAN_APP_KEY

# 各 skill 的百炼 skill_id 映射（能力域目录名 → skill_id）。
# 在百炼控制台「Agent → Skill」中创建/查询后，把 skill_id 填到对应条目。
# 未填写（空字符串）的 skill 只打包、不上传。
$SkillIds = @{
    'hypersku-cli'                             = ''   # TODO: 填写 CLI 能力总览与使用规则 的 skill_id
    'hypersku-auth'                            = ''   # TODO: 填写 CLI 登录与状态管理 的 skill_id
    'hypersku-after-sales'                    = ''   # TODO: 填写 1688售后管理 的 skill_id
    'hypersku-after-sales-apply'              = ''   # TODO: 填写 申请售后 的 skill_id
    'hypersku-customer'                       = ''   # TODO: 填写 客户管理 的 skill_id
    'hypersku-customer-profile-analysis'      = ''   # TODO: 填写 客户画像分析 的 skill_id
    'hypersku-domestic-exception-handling'    = ''   # TODO: 填写 国内异常订单处理 的 skill_id
    'hypersku-domestic-third-trade-exception' = ''   # TODO: 填写 国内第三方交易异常订单管理 的 skill_id
    'hypersku-logistics'                      = ''   # TODO: 填写 物流轨迹查询 的 skill_id
    'hypersku-purchase'                       = ''   # TODO: 填写 采购订单管理 的 skill_id
    'hypersku-warehouse'                      = ''   # TODO: 填写 仓库物流轨迹查询 的 skill_id
}

# 校验环境变量（仅在上传时需要；全部未配置则只打包）
$AnySkillId = $SkillIds.Values | Where-Object { $_ }
$WantUpload = [bool]$AnySkillId
if ($WantUpload -and (-not $DashscopeApiUrl -or -not $DashscopeApiKey)) {
    Write-Host '错误：未配置环境变量 AGENTSTUDIO_URL / BAILIAN_APP_KEY' -ForegroundColor Red
    Write-Host '示例：' -ForegroundColor Yellow
    Write-Host '  $env:AGENTSTUDIO_URL = "https://你的接口地址"' -ForegroundColor Yellow
    Write-Host '  $env:BAILIAN_APP_KEY = "sk-你的key"' -ForegroundColor Yellow
    exit 1
}

Write-Host ''
Write-Host '==============================' -ForegroundColor Cyan
Write-Host '  遍历并压缩所有 skills（多包）' -ForegroundColor Cyan
Write-Host '==============================' -ForegroundColor Cyan

# 创建输出目录（清空旧包，避免残留）
if (Test-Path $OutDir) {
    Get-ChildItem -Path $OutDir -Force | Remove-Item -Recurse -Force
}
New-Item -ItemType Directory -Force -Path $OutDir | Out-Null

# ========== 打包各能力域 skill ==========
# 遍历 skills/ 下一级目录，每个目录（含 SKILL.md）打为一个独立 zip：
# zip 根目录包含 SKILL.md + references/（两级结构，符合平台规范）。
$PackedZips = @()
Get-ChildItem -Path $SkillsSource -Directory | ForEach-Object {
    $skillDir = $_
    $skillName = $skillDir.Name
    $skillFile = Join-Path $skillDir.FullName 'SKILL.md'
    if (-not (Test-Path $skillFile)) {
        Write-Host ("  跳过（缺少 SKILL.md）: " + $skillName) -ForegroundColor DarkGray
        return
    }
    $ZipPath = Join-Path $OutDir ($skillName + '.zip')
    Write-Host ("压缩 skill: " + $skillName + " ...") -ForegroundColor Yellow
    Compress-Archive -Path (Join-Path $skillDir.FullName '*') -DestinationPath $ZipPath -CompressionLevel Optimal -Force
    Write-Host ("  OK 已生成: " + $ZipPath) -ForegroundColor Green
    $PackedZips += [PSCustomObject]@{ Name = $skillName; ZipPath = $ZipPath }
}

# ========== 上传到百炼（逐 skill） ==========
if (-not $WantUpload) {
    Write-Host ''
    Write-Host '跳过上传（$SkillIds 映射表中未配置任何 skill_id），仅完成本地打包。' -ForegroundColor DarkGray
} else {
    foreach ($pkg in $PackedZips) {
        $skillId = $SkillIds[$pkg.Name]
        if (-not $skillId) {
            Write-Host ("  跳过上传（未配置 skill_id）: " + $pkg.Name) -ForegroundColor DarkGray
            continue
        }
        Write-Host ''
        Write-Host ("上传 skill: " + $pkg.Name) -ForegroundColor Yellow

        # 1) 上传 zip 文件 → file_id
        $fileResp = curl.exe -s -X POST "$DashscopeApiUrl/files" `
            -H "Authorization: Bearer $DashscopeApiKey" `
            -F "file=@$($pkg.ZipPath)"
        if ($LASTEXITCODE -ne 0) {
            Write-Host ("  上传文件失败(exit=$LASTEXITCODE): " + $fileResp) -ForegroundColor Red
            continue
        }
        $fileObj = $fileResp | ConvertFrom-Json
        $fileId = $fileObj.id
        if (-not $fileId) {
            Write-Host ("  上传文件失败(未返回 file_id): " + $fileResp) -ForegroundColor Red
            continue
        }
        Write-Host ("  文件ID: " + $fileId) -ForegroundColor Green

        # 2) 上传新版本 → /skills/{skill_id}/versions
        $body = @{ file_id = $fileId } | ConvertTo-Json
        $verResp = curl.exe -s -X POST "$DashscopeApiUrl/skills/$skillId/versions" `
            -H "Authorization: Bearer $DashscopeApiKey" `
            -H "Content-Type: application/json" `
            -d $body
        if ($LASTEXITCODE -ne 0) {
            Write-Host ("  上传版本失败(exit=$LASTEXITCODE): " + $verResp) -ForegroundColor Red
        } else {
            $verObj = $verResp | ConvertFrom-Json
            $status = $verObj.data.status
            Write-Host ("  版本上传成功, status: " + $status) -ForegroundColor Green
        }
    }
}

# ========== 完成 ==========
Write-Host ''
Write-Host ("打包完成 OK（共 " + $PackedZips.Count + " 个独立 skill 包）") -ForegroundColor Green
Write-Host ("  - 输出目录: " + $OutDir)
$PackedZips | ForEach-Object { Write-Host ("  - " + $_.Name + ".zip") }
