@echo off

REM Check if the virtual environment exists
IF NOT EXIST "pyinstaller_venv%suffix%\Scripts\python.exe" (
    REM Create virtual environment
    python -m venv pyinstaller_venv%suffix%
)

REM Activate virtual environment
CALL pyinstaller_venv%suffix%\Scripts\activate

pip install -r requirements.txt
pip install pyinstaller
pip3 install requests

REM Build
python build_pyinstaller.py %*
python -V
REM 检查Python位数
python -c "import struct; print('python is', '32 bit' if struct.calcsize('P') * 8 == 32 else '64 bit')"

REM 退出虚拟环境
CALL deactivate

REM 判断是否是有--onedir参数
set "has_onedir="
for %%i in (%*) do (
    if /i "%%i"=="--onedir" set "has_onedir=1"
)
if defined has_onedir (
    REM 查找并获取dist目录下的exe文件（递归查找）
    for /f "delims=" %%F in ('dir /s /b dist\*%suffix%.exe') do (
        REM 路径中不能包含_internal
        echo %%F | findstr /i "_internal" >nul
        if errorlevel 1 (
            set "exe_file=%%F"
            goto :found_exe
        )
    )
    :found_exe
    REM 获取exe_file的父目录
    for %%I in ("%exe_file%") do (
        set "exe_dir=%%~dpI"
        set "exe_name=%%~nxI"
    )

    REM 移除exe_dir最后的反斜杠
    if "%exe_dir:~-1%" == "\" (
        setlocal enabledelayedexpansion
        set "exe_dir=!exe_dir:~0,-1!"
        endlocal & set "exe_dir=%exe_dir:~0,-1%"
    )

    if not defined exe_dir (
        echo No exe_dir file found in dist directory
        exit /b
    )
    REM 压缩
    echo zip %exe_dir% ....
    python -c "import shutil; shutil.make_archive(r'%exe_dir%', 'zip', r'%exe_dir%')"
    echo zip %exe_dir%.zip done
    REM 计算zip的md5
    echo md5 %exe_dir%.zip ....
    python -c "import hashlib; print(r'%exe_dir%: ' + hashlib.md5(open(r'%exe_dir%.zip', 'rb').read()).hexdigest())" > "%exe_dir%.zip.md5"
    echo md5 done
) else (
    REM 查找并获取dist目录下的exe文件
    for /f "delims=" %%F in ('dir /b dist\*%suffix%.exe') do (
        set "exe_file=dist\%%F"
        set "exe_name=%%F"
    )
    echo %exe_file%

    REM 计算exe的md5
    if defined exe_file (
        python -c "import hashlib; print(r'%exe_name%: ' + hashlib.md5(open(r'%exe_file%', 'rb').read()).hexdigest())" > "%exe_file%.md5"
    ) else (
        echo No exe file found in dist directory
    )
)


