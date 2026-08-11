#Requires -Version 5.1
<#
.SYNOPSIS
    HyperSKU CLI 技能打包脚本

.DESCRIPTION
    将整个 skills/ 目录（统一入口 SKILL.md + 各子域 README.md 与 references/）压缩为一个 zip 包到 build/skills/ 目录。
    zip 根目录直接包含 SKILL.md 及子域目录，不带外层 skills/ 文件夹。

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
$SkillsSource       = Join-Path $ProjectRoot 'skills'
$OutDir             = Join-Path $ProjectRoot 'build\skills'
$DashscopeApiUrl    = $env:AGENTSTUDIO_URL
$DashscopeApiKey    = $env:BAILIAN_APP_KEY

# 统一 skill 的百炼 skill_id（入口已合并为 skills/SKILL.md，整个 skills/ 目录打包为一个 skill）
# 在百炼控制台「Agent → Skill」中创建/查询后，把 skill_id 填到 $SkillId。
# 未填写则只打包、不上传。
$SkillId = 'skill_Zjc1OWZhYWMyNTQzNDU2NDg0Yj'   # TODO: 填写统一 skill 的 skill_id

# 校验环境变量
if (-not $DashscopeApiUrl -or -not $DashscopeApiKey) {
    Write-Host '错误：未配置环境变量 AGENTSTUDIO_URL / BAILIAN_APP_KEY' -ForegroundColor Red
    Write-Host '示例：' -ForegroundColor Yellow
    Write-Host '  $env:AGENTSTUDIO_URL = "https://你的接口地址"' -ForegroundColor Yellow
    Write-Host '  $env:BAILIAN_APP_KEY = "sk-你的key"' -ForegroundColor Yellow
    exit 1
}

Write-Host ''
Write-Host '==============================' -ForegroundColor Cyan
Write-Host '  遍历并压缩所有 skills' -ForegroundColor Cyan
Write-Host '==============================' -ForegroundColor Cyan

# 创建输出目录（清空旧包，避免残留）
if (Test-Path $OutDir) {
    Get-ChildItem -Path $OutDir -Force | Remove-Item -Recurse -Force
}
New-Item -ItemType Directory -Force -Path $OutDir | Out-Null

# ========== 打包统一 skill ==========
# 入口合并后，整个 skills/ 目录视为一个 skill：
# zip 根目录包含 SKILL.md + 各子域目录（README.md + references/），不带外层 skills/ 文件夹。
$ZipPath = Join-Path $OutDir 'hypersku-cli.zip'

Write-Host '压缩统一 skill: hypersku-cli ...' -ForegroundColor Yellow
Compress-Archive -Path (Join-Path $SkillsSource '*') -DestinationPath $ZipPath -CompressionLevel Optimal -Force
Write-Host ("  OK 已生成: " + $ZipPath) -ForegroundColor Green

# ===== 上传到百炼（统一 skill） =====
if (-not $SkillId) {
    Write-Host '  跳过上传（未配置 skill_id）: 请填写 $SkillId' -ForegroundColor DarkGray
} else {
    # 1) 上传 zip 文件 → file_id
    $fileResp = curl.exe -s -X POST "$DashscopeApiUrl/files" `
        -H "Authorization: Bearer $DashscopeApiKey" `
        -F "file=@$ZipPath"
    if ($LASTEXITCODE -ne 0) {
        Write-Host ("  上传文件失败(exit=$LASTEXITCODE): " + $fileResp) -ForegroundColor Red
    } else {
        $fileObj = $fileResp | ConvertFrom-Json
        $fileId = $fileObj.id
        if (-not $fileId) {
            Write-Host ("  上传文件失败(未返回 file_id): " + $fileResp) -ForegroundColor Red
        } else {
            Write-Host ("  文件ID: " + $fileId) -ForegroundColor Green

            # 2) 上传新版本 → /skills/{skill_id}/versions
            $body = @{ file_id = $fileId } | ConvertTo-Json
            $verResp = curl.exe -s -X POST "$DashscopeApiUrl/skills/$SkillId/versions" `
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
}

# ========== 完成 ==========
Write-Host ''
Write-Host '打包完成 OK（统一 skill: hypersku-cli）' -ForegroundColor Green
Write-Host ("  - 输出目录: " + $OutDir)
