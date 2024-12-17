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
