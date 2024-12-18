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
python build_pyinstaller.py --onedir  %*
python -V
REM 检查Python位数
python -c "import struct; print('python is', '32 bit' if struct.calcsize('P') * 8 == 32 else '64 bit')"

REM 退出虚拟环境
CALL deactivate

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
)
set "exe_dir=%exe_dir:~0,-1%"


REM 压缩
echo zip %exe_dir% ....
python -c "import shutil; shutil.make_archive(r'%exe_dir%', 'zip', r'%exe_dir%')"
echo zip %exe_dir%.zip done
REM 计算zip的md5
echo md5 %exe_dir%.zip ....
python -c "import hashlib; print(hashlib.md5(open(r'%exe_dir%.zip', 'rb').read()).hexdigest())" > "%exe_dir%.zip.md5"
echo md5 done




