REM Check if the virtual environment exists
IF NOT EXIST "pyinstaller_venv" (
    REM Create virtual environment
    python -m venv pyinstaller_venv
)

REM Activate virtual environment
CALL pyinstaller_venv\Scripts\activate

pip install -r requirements.txt
pip install pyinstaller

REM Build
python build_pyinstaller.py %*

REM 查找并获取dist目录下的exe文件
for /f "delims=" %%F in ('dir /b dist\*.exe') do (
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
